package targets

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mitchellh/mapstructure"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	txmgrcommon "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/forwarder"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	evm "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	relayevmtypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

var forwardABI = evmtypes.MustGetABI(forwarder.KeystoneForwarderMetaData.ABI)

func InitializeWrite(registry core.CapabilitiesRegistry, legacyEVMChains legacyevm.LegacyChainContainer, lggr logger.Logger) error {
	for _, chain := range legacyEVMChains.Slice() {
		capability := NewEvmWrite(chain, lggr)
		if err := registry.Add(context.TODO(), capability); err != nil {
			return err
		}
	}
	return nil
}

var (
	_ capabilities.ActionCapability = &EvmWrite{}
)

const defaultGasLimit = 200000

type EvmWrite struct {
	chain legacyevm.Chain
	capabilities.CapabilityInfo
	lggr logger.Logger
}

func NewEvmWrite(chain legacyevm.Chain, lggr logger.Logger) *EvmWrite {
	// generate ID based on chain selector
	name := fmt.Sprintf("write_%v", chain.ID())
	chainName, err := chainselectors.NameFromChainId(chain.ID().Uint64())
	if err == nil {
		name = fmt.Sprintf("write_%v", chainName)
	}

	info := capabilities.MustNewCapabilityInfo(
		name,
		capabilities.CapabilityTypeTarget,
		"Write target.",
		"v1.0.0",
		nil,
	)

	return &EvmWrite{
		chain,
		info,
		lggr.Named("EvmWrite"),
	}
}

type EvmConfig struct {
	ChainID uint
	Address string
}

// TODO: enforce required key presence

func parseConfig(rawConfig *values.Map) (EvmConfig, error) {
	var config EvmConfig
	configAny, err := rawConfig.Unwrap()
	if err != nil {
		return config, err
	}
	err = mapstructure.Decode(configAny, &config)
	return config, err
}

func success() <-chan capabilities.CapabilityResponse {
	callback := make(chan capabilities.CapabilityResponse)
	go func() {
		// TODO: cast tx.Error to Err (or Value to Value?)
		callback <- capabilities.CapabilityResponse{
			Value: nil,
			Err:   nil,
		}
		close(callback)
	}()
	return callback
}

func (cap *EvmWrite) Execute(ctx context.Context, request capabilities.CapabilityRequest) (<-chan capabilities.CapabilityResponse, error) {
	cap.lggr.Debugw("Execute", "request", request)

	txm := cap.chain.TxManager()

	config := cap.chain.Config().EVM().ChainWriter()

	reqConfig, err := parseConfig(request.Config)
	if err != nil {
		return nil, err
	}
	var inputs struct {
		Report     []byte
		Signatures [][]byte
	}
	if err = request.Inputs.UnwrapTo(&inputs); err != nil {
		return nil, err
	}

	if inputs.Report == nil {
		// We received any empty report -- this means we should skip transmission.
		cap.lggr.Debugw("Skipping empty report", "request", request)
		return success(), nil
	}

	// TODO: validate encoded report is prefixed with workflowID and executionID that match the request meta

	// Check whether value was already transmitted on chain
	cr, err := evm.NewChainReaderService(ctx, cap.lggr, cap.chain.LogPoller(), cap.chain.Client(), relayevmtypes.ChainReaderConfig{
		Contracts: map[string]relayevmtypes.ChainContractReader{
			"forwarder": {
				ContractABI: forwarder.KeystoneForwarderABI,
				Configs: map[string]*relayevmtypes.ChainReaderDefinition{
					"getTransmitter": {
						ChainSpecificName: "getTransmitter",
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	var transmitter common.Address
	cr.Bind(ctx, []commontypes.BoundContract{{
		Address: config.ForwarderAddress().String(),
		Name:    "forwarder",
		// Pending: false, // ???
	}})
	queryInputs := struct {
		Receiver            string
		WorkflowExecutionID []byte
	}{
		Receiver:            reqConfig.Address,
		WorkflowExecutionID: []byte(request.Metadata.WorkflowExecutionID),
	}
	if err := cr.GetLatestValue(ctx, "forwarder", "getTransmitter", queryInputs, &transmitter); err != nil {
		return nil, err
	}
	if transmitter != common.HexToAddress("0x0") {
		// report already transmitted, early return
		return success(), nil
	}

	// construct forwarder payload
	calldata, err := forwardABI.Pack("report", common.HexToAddress(reqConfig.Address), inputs.Report, inputs.Signatures)
	if err != nil {
		return nil, err
	}

	txMeta := &txmgr.TxMeta{
		// FwdrDestAddress could also be set for better logging but it's used for various purposes around Operator Forwarders
		WorkflowExecutionID: &request.Metadata.WorkflowExecutionID,
	}
	req := txmgr.TxRequest{
		FromAddress:    config.FromAddress().Address(),
		ToAddress:      config.ForwarderAddress().Address(),
		EncodedPayload: calldata,
		FeeLimit:       uint64(defaultGasLimit),
		Meta:           txMeta,
		Strategy:       txmgrcommon.NewSendEveryStrategy(),
		Checker: txmgr.TransmitCheckerSpec{
			CheckerType: txmgr.TransmitCheckerTypeSimulate,
		},
		// SignalCallback:   true, TODO: add code that checks if a workflow id is present, if so, route callback to chainwriter rather than pipeline
	}
	tx, err := txm.CreateTransaction(ctx, req)
	if err != nil {
		return nil, err
	}
	cap.lggr.Debugw("Transaction submitted", "request", request, "transaction", tx)
	return success(), nil
}

func (cap *EvmWrite) RegisterToWorkflow(ctx context.Context, request capabilities.RegisterToWorkflowRequest) error {
	return nil
}

func (cap *EvmWrite) UnregisterFromWorkflow(ctx context.Context, request capabilities.UnregisterFromWorkflowRequest) error {
	return nil
}
