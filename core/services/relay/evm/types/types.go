package types

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"gopkg.in/guregu/null.v2"

	"github.com/smartcontractkit/chainlink-common/pkg/codec"
	"github.com/smartcontractkit/chainlink-common/pkg/services"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
)

type ChainReaderConfig struct {
	// ChainContractReaders key is contract name
	ChainContractReaders map[string]ChainContractReader `json:"chainContractReaders"`
}

type CodecConfig struct {
	// ChainCodecConfigs is the type's name for the codec
	ChainCodecConfigs map[string]ChainCodedConfig `json:"chainCodecConfig"`
}

type ChainCodedConfig struct {
	TypeAbi         string `json:"typeAbi"`
	ModifierConfigs codec.ModifiersConfig
}

type ChainContractReader struct {
	ContractABI string `json:"contractABI"`
	// key is genericName from config
	ChainReaderDefinitions map[string]ChainReaderDefinition `json:"chainReaderDefinitions"`
}

type ChainReaderDefinition struct {
	ChainSpecificName   string `json:"chainSpecificName"` // chain specific contract method name or event type.
	CacheEnabled        bool   `json:"cacheEnabled"`
	ReadType            `json:"readType"`
	InputModifications  codec.ModifiersConfig `json:"input_modifications"`
	OutputModifications codec.ModifiersConfig `json:"output_modifications"`
}

type ReadType int64

const (
	Method ReadType = 0
	Event  ReadType = 1
)

type RelayConfig struct {
	ChainID                *big.Big           `json:"chainID"`
	FromBlock              uint64             `json:"fromBlock"`
	EffectiveTransmitterID null.String        `json:"effectiveTransmitterID"`
	ConfigContractAddress  *common.Address    `json:"configContractAddress"`
	ChainReader            *ChainReaderConfig `json:"chainReader"`
	Codec                  *CodecConfig       `json:"codec"`

	// Contract-specific
	SendingKeys pq.StringArray `json:"sendingKeys"`

	// Mercury-specific
	FeedID *common.Hash `json:"feedID"`
}

var ErrBadRelayConfig = errors.New("bad relay config")

type RelayOpts struct {
	// TODO BCF-2508 -- should anyone ever get the raw config bytes that are embedded in args? if not,
	// make this private and wrap the arg fields with funcs on RelayOpts
	commontypes.RelayArgs
	c *RelayConfig
}

func NewRelayOpts(args commontypes.RelayArgs) *RelayOpts {
	return &RelayOpts{
		RelayArgs: args,
		c:         nil, // lazy initialization
	}
}

func (o *RelayOpts) RelayConfig() (RelayConfig, error) {
	var empty RelayConfig
	//TODO this should be done once and the error should be cached
	if o.c == nil {
		var c RelayConfig
		err := json.Unmarshal(o.RelayArgs.RelayConfig, &c)
		if err != nil {
			return empty, fmt.Errorf("%w: failed to deserialize relay config: %w", ErrBadRelayConfig, err)
		}
		o.c = &c
	}
	return *o.c, nil
}

type ConfigPoller interface {
	ocrtypes.ContractConfigTracker

	Start()
	Close() error
	Replay(ctx context.Context, fromBlock int64) error
}

// TODO(FUN-668): Migrate this fully into commontypes.FunctionsProvider
type FunctionsProvider interface {
	commontypes.FunctionsProvider
	LogPollerWrapper() LogPollerWrapper
}

type OracleRequest struct {
	RequestId           [32]byte
	RequestingContract  common.Address
	RequestInitiator    common.Address
	SubscriptionId      uint64
	SubscriptionOwner   common.Address
	Data                []byte
	DataVersion         uint16
	Flags               [32]byte
	CallbackGasLimit    uint64
	TxHash              common.Hash
	CoordinatorContract common.Address
	OnchainMetadata     []byte
}

type OracleResponse struct {
	RequestId [32]byte
}

type RouteUpdateSubscriber interface {
	UpdateRoutes(activeCoordinator common.Address, proposedCoordinator common.Address) error
}

// A LogPoller wrapper that understands router proxy contracts
//
//go:generate mockery --quiet --name LogPollerWrapper --output ./mocks/ --case=underscore
type LogPollerWrapper interface {
	services.Service
	LatestEvents() ([]OracleRequest, []OracleResponse, error)

	// TODO (FUN-668): Remove from the LOOP interface and only use internally within the EVM relayer
	SubscribeToUpdates(name string, subscriber RouteUpdateSubscriber)
}
