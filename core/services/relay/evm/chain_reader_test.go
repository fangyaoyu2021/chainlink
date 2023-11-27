package evm

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	mocklogpoller "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

type chainReaderTestHelper struct {
}

func (crTestHelper chainReaderTestHelper) makeChainReaderConfig(abi string, params map[string]any, retValues []string) evmtypes.ChainReaderConfig {
	return evmtypes.ChainReaderConfig{
		ChainContractReaders: map[string]evmtypes.ChainContractReader{
			"MyContract": {
				ContractABI: abi,
				ChainReaderDefinitions: map[string]evmtypes.ChainReaderDefinition{
					"MyGenericMethod": {
						ChainSpecificName: "name",
						Params:            params,
						ReturnValues:      retValues,
						CacheEnabled:      false,
						ReadType:          evmtypes.Method,
					},
				},
			},
		},
	}
}

func (crTestHelper chainReaderTestHelper) makeChainReaderConfigFromStrings(abi string, chainReadingDefinitions string) (evmtypes.ChainReaderConfig, error) {
	chainReaderConfigTemplate := `{
	   "chainContractReaders": {
	     "testContract": {
			   "contractName": "testContract",
			   "contractABI":  "[%s]",
			   "chainReaderDefinitions": {
					%s
				}
	     }
	   }
	}`

	abi = strings.Replace(abi, `"`, `\"`, -1)
	formattedCfgJsonString := fmt.Sprintf(chainReaderConfigTemplate, abi, chainReadingDefinitions)
	var chainReaderConfig evmtypes.ChainReaderConfig
	err := json.Unmarshal([]byte(formattedCfgJsonString), &chainReaderConfig)
	return chainReaderConfig, err
}

func TestNewChainReader(t *testing.T) {
	lggr := logger.TestLogger(t)
	lp := mocklogpoller.NewLogPoller(t)
	chain := mocks.NewChain(t)
	contractID := testutils.NewAddress()
	contractABI := `[{"inputs":[],"name":"name","outputs":[{"internalType":"string","name":"result","type":"string"}],"stateMutability":"view","type":"function"}]`

	t.Run("happy path", func(t *testing.T) {
		chainReaderConfig := chainReaderTestHelper{}.makeChainReaderConfig(contractABI, map[string]any{}, []string{"result"})
		chain.On("LogPoller").Return(lp)
		_, err := NewChainReaderService(lggr, chain.LogPoller(), contractID, chainReaderConfig)
		assert.NoError(t, err)
	})

	t.Run("invalid config", func(t *testing.T) {
		invalidChainReaderConfig := chainReaderTestHelper{}.makeChainReaderConfig(contractABI, map[string]any{}, []string{"result", "extraResult"}) // 2 results required but abi includes only one
		_, err := NewChainReaderService(lggr, chain.LogPoller(), contractID, invalidChainReaderConfig)
		assert.ErrorIs(t, err, commontypes.ErrInvalidConfig)
		assert.ErrorContains(t, err, "return values: [result,extraResult] don't match abi method outputs: [result]")
	})

	t.Run("ChainReader config is empty", func(t *testing.T) {
		emptyChainReaderConfig := evmtypes.ChainReaderConfig{}
		_, err := NewChainReaderService(lggr, chain.LogPoller(), contractID, emptyChainReaderConfig)
		assert.ErrorIs(t, err, commontypes.ErrInvalidConfig)
		assert.ErrorContains(t, err, "config is empty")
	})
}

func TestChainReaderStartClose(t *testing.T) {
	lggr := logger.TestLogger(t)
	lp := mocklogpoller.NewLogPoller(t)
	cr := chainReader{
		lggr: lggr,
		lp:   lp,
	}
	err := cr.Start(testutils.Context(t))
	assert.NoError(t, err)
	err = cr.Close()
	assert.NoError(t, err)
}

func TestValidateChainReaderConfig_HappyPath(t *testing.T) {
	type testCase struct {
		name                    string
		abiInput                string
		chainReadingDefinitions string
	}

	var testCases []testCase
	testCases = append(testCases, testCase{
		name:     "median abi",
		abiInput: `{"inputs":[{"internalType":"uint32","name":"_maximumGasPrice","type":"uint32"},{"internalType":"uint32","name":"_reasonableGasPrice","type":"uint32"},{"internalType":"uint32","name":"_microLinkPerEth","type":"uint32"},{"internalType":"uint32","name":"_linkGweiPerObservation","type":"uint32"},{"internalType":"uint32","name":"_linkGweiPerTransmission","type":"uint32"},{"internalType":"address","name":"_link","type":"address"},{"internalType":"address","name":"_validator","type":"address"},{"internalType":"int192","name":"_minAnswer","type":"int192"},{"internalType":"int192","name":"_maxAnswer","type":"int192"},{"internalType":"contractAccessControllerInterface","name":"_billingAccessController","type":"address"},{"internalType":"contractAccessControllerInterface","name":"_requesterAccessController","type":"address"},{"internalType":"uint8","name":"_decimals","type":"uint8"},{"internalType":"string","name":"_description","type":"string"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"int256","name":"current","type":"int256"},{"indexed":true,"internalType":"uint256","name":"roundId","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"updatedAt","type":"uint256"}],"name":"AnswerUpdated","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"contractAccessControllerInterface","name":"old","type":"address"},{"indexed":false,"internalType":"contractAccessControllerInterface","name":"current","type":"address"}],"name":"BillingAccessControllerSet","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint32","name":"maximumGasPrice","type":"uint32"},{"indexed":false,"internalType":"uint32","name":"reasonableGasPrice","type":"uint32"},{"indexed":false,"internalType":"uint32","name":"microLinkPerEth","type":"uint32"},{"indexed":false,"internalType":"uint32","name":"linkGweiPerObservation","type":"uint32"},{"indexed":false,"internalType":"uint32","name":"linkGweiPerTransmission","type":"uint32"}],"name":"BillingSet","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint32","name":"previousConfigBlockNumber","type":"uint32"},{"indexed":false,"internalType":"uint64","name":"configCount","type":"uint64"},{"indexed":false,"internalType":"address[]","name":"signers","type":"address[]"},{"indexed":false,"internalType":"address[]","name":"transmitters","type":"address[]"},{"indexed":false,"internalType":"uint8","name":"threshold","type":"uint8"},{"indexed":false,"internalType":"uint64","name":"encodedConfigVersion","type":"uint64"},{"indexed":false,"internalType":"bytes","name":"encoded","type":"bytes"}],"name":"ConfigSet","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint256","name":"roundId","type":"uint256"},{"indexed":true,"internalType":"address","name":"startedBy","type":"address"},{"indexed":false,"internalType":"uint256","name":"startedAt","type":"uint256"}],"name":"NewRound","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint32","name":"aggregatorRoundId","type":"uint32"},{"indexed":false,"internalType":"int192","name":"answer","type":"int192"},{"indexed":false,"internalType":"address","name":"transmitter","type":"address"},{"indexed":false,"internalType":"int192[]","name":"observations","type":"int192[]"},{"indexed":false,"internalType":"bytes","name":"observers","type":"bytes"},{"indexed":false,"internalType":"bytes32","name":"rawReportContext","type":"bytes32"}],"name":"NewTransmission","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"transmitter","type":"address"},{"indexed":false,"internalType":"address","name":"payee","type":"address"},{"indexed":false,"internalType":"uint256","name":"amount","type":"uint256"}],"name":"OraclePaid","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":true,"internalType":"address","name":"to","type":"address"}],"name":"OwnershipTransferRequested","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":true,"internalType":"address","name":"to","type":"address"}],"name":"OwnershipTransferred","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"transmitter","type":"address"},{"indexed":true,"internalType":"address","name":"current","type":"address"},{"indexed":true,"internalType":"address","name":"proposed","type":"address"}],"name":"PayeeshipTransferRequested","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"transmitter","type":"address"},{"indexed":true,"internalType":"address","name":"previous","type":"address"},{"indexed":true,"internalType":"address","name":"current","type":"address"}],"name":"PayeeshipTransferred","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"contractAccessControllerInterface","name":"old","type":"address"},{"indexed":false,"internalType":"contractAccessControllerInterface","name":"current","type":"address"}],"name":"RequesterAccessControllerSet","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"requester","type":"address"},{"indexed":false,"internalType":"bytes16","name":"configDigest","type":"bytes16"},{"indexed":false,"internalType":"uint32","name":"epoch","type":"uint32"},{"indexed":false,"internalType":"uint8","name":"round","type":"uint8"}],"name":"RoundRequested","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"previous","type":"address"},{"indexed":true,"internalType":"address","name":"current","type":"address"}],"name":"ValidatorUpdated","type":"event"},{"inputs":[],"name":"LINK","outputs":[{"internalType":"contractLinkTokenInterface","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"acceptOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"_transmitter","type":"address"}],"name":"acceptPayeeship","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"billingAccessController","outputs":[{"internalType":"contractAccessControllerInterface","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"decimals","outputs":[{"internalType":"uint8","name":"","type":"uint8"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"description","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"_roundId","type":"uint256"}],"name":"getAnswer","outputs":[{"internalType":"int256","name":"","type":"int256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getBilling","outputs":[{"internalType":"uint32","name":"maximumGasPrice","type":"uint32"},{"internalType":"uint32","name":"reasonableGasPrice","type":"uint32"},{"internalType":"uint32","name":"microLinkPerEth","type":"uint32"},{"internalType":"uint32","name":"linkGweiPerObservation","type":"uint32"},{"internalType":"uint32","name":"linkGweiPerTransmission","type":"uint32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint80","name":"_roundId","type":"uint80"}],"name":"getRoundData","outputs":[{"internalType":"uint80","name":"roundId","type":"uint80"},{"internalType":"int256","name":"answer","type":"int256"},{"internalType":"uint256","name":"startedAt","type":"uint256"},{"internalType":"uint256","name":"updatedAt","type":"uint256"},{"internalType":"uint80","name":"answeredInRound","type":"uint80"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"_roundId","type":"uint256"}],"name":"getTimestamp","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"latestAnswer","outputs":[{"internalType":"int256","name":"","type":"int256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"latestConfigDetails","outputs":[{"internalType":"uint32","name":"configCount","type":"uint32"},{"internalType":"uint32","name":"blockNumber","type":"uint32"},{"internalType":"bytes16","name":"configDigest","type":"bytes16"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"latestRound","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"latestRoundData","outputs":[{"internalType":"uint80","name":"roundId","type":"uint80"},{"internalType":"int256","name":"answer","type":"int256"},{"internalType":"uint256","name":"startedAt","type":"uint256"},{"internalType":"uint256","name":"updatedAt","type":"uint256"},{"internalType":"uint80","name":"answeredInRound","type":"uint80"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"latestTimestamp","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"latestTransmissionDetails","outputs":[{"internalType":"bytes16","name":"configDigest","type":"bytes16"},{"internalType":"uint32","name":"epoch","type":"uint32"},{"internalType":"uint8","name":"round","type":"uint8"},{"internalType":"int192","name":"latestAnswer","type":"int192"},{"internalType":"uint64","name":"latestTimestamp","type":"uint64"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"linkAvailableForPayment","outputs":[{"internalType":"int256","name":"availableBalance","type":"int256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"maxAnswer","outputs":[{"internalType":"int192","name":"","type":"int192"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"minAnswer","outputs":[{"internalType":"int192","name":"","type":"int192"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_signerOrTransmitter","type":"address"}],"name":"oracleObservationCount","outputs":[{"internalType":"uint16","name":"","type":"uint16"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_transmitter","type":"address"}],"name":"owedPayment","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"owner","outputs":[{"internalType":"addresspayable","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"requestNewRound","outputs":[{"internalType":"uint80","name":"","type":"uint80"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"requesterAccessController","outputs":[{"internalType":"contractAccessControllerInterface","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint32","name":"_maximumGasPrice","type":"uint32"},{"internalType":"uint32","name":"_reasonableGasPrice","type":"uint32"},{"internalType":"uint32","name":"_microLinkPerEth","type":"uint32"},{"internalType":"uint32","name":"_linkGweiPerObservation","type":"uint32"},{"internalType":"uint32","name":"_linkGweiPerTransmission","type":"uint32"}],"name":"setBilling","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"contractAccessControllerInterface","name":"_billingAccessController","type":"address"}],"name":"setBillingAccessController","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address[]","name":"_signers","type":"address[]"},{"internalType":"address[]","name":"_transmitters","type":"address[]"},{"internalType":"uint8","name":"_threshold","type":"uint8"},{"internalType":"uint64","name":"_encodedConfigVersion","type":"uint64"},{"internalType":"bytes","name":"_encoded","type":"bytes"}],"name":"setConfig","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address[]","name":"_transmitters","type":"address[]"},{"internalType":"address[]","name":"_payees","type":"address[]"}],"name":"setPayees","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"contractAccessControllerInterface","name":"_requesterAccessController","type":"address"}],"name":"setRequesterAccessController","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"_newValidator","type":"address"}],"name":"setValidator","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"_to","type":"address"}],"name":"transferOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"_transmitter","type":"address"},{"internalType":"address","name":"_proposed","type":"address"}],"name":"transferPayeeship","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes","name":"_report","type":"bytes"},{"internalType":"bytes32[]","name":"_rs","type":"bytes32[]"},{"internalType":"bytes32[]","name":"_ss","type":"bytes32[]"},{"internalType":"bytes32","name":"_rawVs","type":"bytes32"}],"name":"transmit","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"transmitters","outputs":[{"internalType":"address[]","name":"","type":"address[]"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"validator","outputs":[{"internalType":"contractAggregatorValidatorInterface","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"version","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_recipient","type":"address"},{"internalType":"uint256","name":"_amount","type":"uint256"}],"name":"withdrawFunds","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"_transmitter","type":"address"}],"name":"withdrawPayment","outputs":[],"stateMutability":"nonpayable","type":"function"}`,
		chainReadingDefinitions: ` "latestTransmissionDetails":{
											"chainSpecificName": "latestTransmissionDetails",
											"returnValues": [
												"configDigest",
												"epoch",
												"round",
												"latestAnswer",
												"latestTimestamp"
											],
											"readType": 0
											},
									"latestRoundRequested":{
											"chainSpecificName": "RoundRequested",
											"params": {
												"requester": ""
											},
											"returnValues": [
												"requester",
												"configDigest",
												"epoch",
												"round"
											],
											"readType": 1
										}`,
	})
	testCases = append(testCases,
		testCase{
			name:     "eventWithNoIndexedTopics",
			abiInput: `{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint112","name":"reserve0","type":"uint112"},{"indexed":false,"internalType":"uint112","name":"reserve1","type":"uint112"}],"name":"Sync","type":"event"}`,
			chainReadingDefinitions: ` "Sync":{
											"chainSpecificName": "Sync",
											"returnValues": [
												"reserve0",
												"reserve1"
											],
											"readType": 1
										}`,
		})

	testCases = append(testCases,
		testCase{
			name:     "eventWithMultipleIndexedTopics",
			abiInput: `{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"sender","type":"address"},{"indexed":false,"internalType":"uint256","name":"amount0In","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"amount1In","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"amount0Out","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"amount1Out","type":"uint256"},{"indexed":true,"internalType":"address","name":"to","type":"address"}],"name":"Swap","type":"event"}`,
			chainReadingDefinitions: `"Swap":{
											"chainSpecificName": "Swap",
											"params":{
												"sender": "0x0",
												"to": "0x0"
											},
											"returnValues": [
												"sender",
												"amount0In",
												"amount1In",
												"amount0Out",
												"amount1Out",
												"to"
											],
											"readType": 1
										}`,
		})

	testCases = append(testCases,
		testCase{
			name:     "methodWithOneParamAndMultipleResponses",
			abiInput: `{"constant":true,"inputs":[{"internalType":"address","name":"_user","type":"address"}],"name":"getUserAccountData","outputs":[{"internalType":"uint256","name":"totalLiquidityETH","type":"uint256"},{"internalType":"uint256","name":"totalCollateralETH","type":"uint256"},{"internalType":"uint256","name":"totalBorrowsETH","type":"uint256"},{"internalType":"uint256","name":"totalFeesETH","type":"uint256"},{"internalType":"uint256","name":"availableBorrowsETH","type":"uint256"},{"internalType":"uint256","name":"currentLiquidationThreshold","type":"uint256"},{"internalType":"uint256","name":"ltv","type":"uint256"},{"internalType":"uint256","name":"healthFactor","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"}`,
			chainReadingDefinitions: `"getUserAccountData":{
											"chainSpecificName": "getUserAccountData",
											"params":{
												"_user": "0x0"
											},
											"returnValues": [
												"totalLiquidityETH",
												"totalCollateralETH",
												"totalBorrowsETH",
												"totalFeesETH",
												"availableBorrowsETH",
												"currentLiquidationThreshold",
												"healthFactor"
											],
											"readType": 0
										}`,
		})

	testCases = append(testCases,
		testCase{
			name:     "methodWithNoParamsAndOneResponseWithNoName",
			abiInput: `{"inputs":[],"name":"name","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"}`,
			chainReadingDefinitions: `"name":{
											"chainSpecificName": "name",
											"returnValues": [
												""
											],
											"readType": 0
										}`,
		})

	testCases = append(testCases,
		testCase{
			name:     "methodWithMultipleParamsAndOneResult",
			abiInput: `{"inputs":[{"internalType":"address","name":"_input","type":"address"},{"internalType":"address","name":"_output","type":"address"},{"internalType":"uint256","name":"_inputQuantity","type":"uint256"}],"name":"getSwapOutput","outputs":[{"internalType":"uint256","name":"swapOutput","type":"uint256"}],"stateMutability":"view","type":"function"}`,
			chainReadingDefinitions: `"getSwapOutput":{
											"chainSpecificName": "getSwapOutput",
											"params":{
												"_input":"0x0",
												"_output":"0x0",
												"_inputQuantity":"0x0"
											},
											"returnValues": [
												"swapOutput"
											],
											"readType": 0
										}`,
		})

	// TODO BCF-2789 how to handle return values for tuples
	/*testCases = append(testCases,
	testCase{
		name: "methodWithOneParamAndMultipleTupleResponse",
		 struct BassetPersonal {
		    // Address of the bAsset
		    address addr;
		    // Address of the bAsset
		    address integrator;
		    // An ERC20 can charge transfer fee, for example USDT, DGX tokens.
		    bool hasTxFee; // takes a byte in storage
		    // Status of the bAsset
		    BassetStatus status;
		}

		// Status of the Basset - has it broken its peg?
		enum BassetStatus {
		    Default,
		    Normal,
		    BrokenBelowPeg,
		    BrokenAbovePeg,
		    Blacklisted,
		    Liquidating,
		    Liquidated,
		    Failed
		}

		struct BassetData {
		    // 1 Basset * ratio / ratioScale == x Masset (relative value)
		    // If ratio == 10e8 then 1 bAsset = 10 mAssets
		    // A ratio is divised as 10^(18-tokenDecimals) * measurementMultiple(relative value of 1 base unit)
		    uint128 ratio;
		    // Amount of the Basset that is held in Collateral
		    uint128 vaultBalance;
		}
		abiInput: `{"inputs":[{"internalType":"address","name":"_bAsset","type":"address"}],"name":"getBasset","outputs":[{"components":[{"internalType":"address","name":"addr","type":"address"},{"internalType":"address","name":"integrator","type":"address"},{"internalType":"bool","name":"hasTxFee","type":"bool"},{"internalType":"enum BassetStatus","name":"status","type":"uint8"}],"internalType":"struct BassetPersonal","name":"personal","type":"tuple"},{"components":[{"internalType":"uint128","name":"ratio","type":"uint128"},{"internalType":"uint128","name":"vaultBalance","type":"uint128"}],"internalType":"struct BassetData","name":"vaultData","type":"tuple"}],"stateMutability":"view","type":"function"}`,
		chainReadingDefinitions: `getBasset:{
										chainSpecificName: getBasset,
										params:{
											_bAsset:"0x0",
										},
										returnValues: [
											TODO,
										]
										readType: 0,
									}`,
	})
	*/

	// TODO how to handle return values for tuples
	/*
		testCases = append(testCases,
			testCase{
				name: "methodWithNoParamsAndTupleResponse",
				 struct FeederConfig {
					uint256 supply;
					uint256 a;
					WeightLimits limits;
				}

				struct WeightLimits {
					uint128 min;
					uint128 max;
				}
				abiInput: `{"inputs":[],"name":"getConfig","outputs":[{"components":[{"internalType":"uint256","name":"supply","type":"uint256"},{"internalType":"uint256","name":"a","type":"uint256"},{"components":[{"internalType":"uint128","name":"min","type":"uint128"},{"internalType":"uint128","name":"max","type":"uint128"}],"internalType":"struct WeightLimits","name":"limits","type":"tuple"}],"internalType":"struct FeederConfig","name":"config","type":"tuple"}],"stateMutability":"view","type":"function"}`,
				chainReadingDefinitions: `getConfig:{
												chainSpecificName: getConfig,
												params:{},
												returnValues: [
													TODO,
												]
												readType: 0,
											}`,
			})*/

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg, err := chainReaderTestHelper{}.makeChainReaderConfigFromStrings(tc.abiInput, tc.chainReadingDefinitions)
			assert.NoError(t, err)
			assert.NoError(t, validateChainReaderConfig(cfg))
		})
	}

	t.Run("large config with all test cases", func(t *testing.T) {
		var largeABI string
		var manyChainReadingDefinitions string
		for _, tc := range testCases {
			largeABI += tc.abiInput + ","
			manyChainReadingDefinitions += tc.chainReadingDefinitions + ","
		}

		largeABI = largeABI[:len(largeABI)-1]
		manyChainReadingDefinitions = manyChainReadingDefinitions[:len(manyChainReadingDefinitions)-1]
		cfg, err := chainReaderTestHelper{}.makeChainReaderConfigFromStrings(largeABI, manyChainReadingDefinitions)
		assert.NoError(t, err)
		assert.NoError(t, validateChainReaderConfig(cfg))
	})
}

func TestValidateChainReaderConfig_BadPath(t *testing.T) {
	type testCase struct {
		name                    string
		abiInput                string
		chainReadingDefinitions string
		expected                error
	}

	var testCases []testCase
	mismatchedEventArgumentsTestABI := `{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"sender","type":"address"}],"name":"Swap","type":"event"}`
	testCases = append(testCases,
		testCase{
			name:     "mismatched abi and event chain reading param values",
			abiInput: mismatchedEventArgumentsTestABI,
			chainReadingDefinitions: `"Swap":{
													"chainSpecificName": "Swap",
													"params":{
														"malformedParam": "0x0"
													},
													"returnValues": [
														"sender"
													],
													"readType": 1
												}`,
			expected: fmt.Errorf("invalid chain reading definition: \"Swap\" for contract: \"testContract\", err: params: [malformedParam] don't match abi event indexed inputs: [sender]"),
		})

	testCases = append(testCases,
		testCase{
			name:     "mismatched abi and event chain reading return values",
			abiInput: mismatchedEventArgumentsTestABI,
			chainReadingDefinitions: `"Swap":{
													"chainSpecificName": "Swap",
													"params":{
														"sender": "0x0"
													},
													"returnValues": [
														"malformedParam"
													],
													"readType": 1
												}`,
			expected: fmt.Errorf("invalid chain reading definition: \"Swap\" for contract: \"testContract\", err: return values: [malformedParam] don't match abi event inputs: [sender]"),
		})

	mismatchedFunctionArgumentsTestABI := `{"constant":true,"inputs":[{"internalType":"address","name":"from","type":"address"}],"name":"Swap","outputs":[{"internalType":"address","name":"to","type":"address"}],"payable":false,"stateMutability":"view","type":"function"}`
	testCases = append(testCases,
		testCase{
			name:     "mismatched abi and method chain reading param values",
			abiInput: mismatchedFunctionArgumentsTestABI,
			chainReadingDefinitions: `"Swap":{
											"chainSpecificName": "Swap",
											"params":{
												"malformedParam": "0x0"
											},
											"returnValues": [
												"to"
											],
											"readType": 0
										}`,
			expected: fmt.Errorf("invalid chain reading definition: \"Swap\" for contract: \"testContract\", err: params: [malformedParam] don't match abi method inputs: [from]"),
		},
	)

	testCases = append(testCases,
		testCase{
			name:     "mismatched abi and method chain reading return values",
			abiInput: mismatchedFunctionArgumentsTestABI,
			chainReadingDefinitions: `"Swap":{
											"chainSpecificName": "Swap",
											"params":{
												"from": "0x0"
											},
											"returnValues": [
												"malformedReturnValue"
											],
											"readType": 0
										}`,
			expected: fmt.Errorf("invalid chain reading definition: \"Swap\" for contract: \"testContract\", err: return values: [malformedReturnValue] don't match abi method outputs: [to]"),
		},
	)

	testCases = append(testCases,
		testCase{
			name:     "event doesn't exist",
			abiInput: `{"constant":true,"inputs":[],"name":"someName","outputs":[],"payable":false,"stateMutability":"view","type":"function"}`,
			chainReadingDefinitions: `"TestMethod":{
											"chainSpecificName": "Swap",
											"readType": 1
										}`,
			expected: fmt.Errorf("invalid chain reading definition: \"TestMethod\" for contract: \"testContract\", err: event: Swap doesn't exist"),
		},
	)

	testCases = append(testCases,
		testCase{
			name:     "method doesn't exist",
			abiInput: `{"constant":true,"inputs":[],"name":"someName","outputs":[],"payable":false,"stateMutability":"view","type":"function"}`,
			chainReadingDefinitions: `"TestMethod":{
											"chainSpecificName": "Swap",
											"readType": 0
										}`,
			expected: fmt.Errorf("invalid chain reading definition: \"TestMethod\" for contract: \"testContract\", err: method: \"Swap\" doesn't exist"),
		},
	)

	testCases = append(testCases, testCase{
		name:     "invalid abi",
		abiInput: `broken abi`,
		chainReadingDefinitions: `"TestMethod":{
											"chainSpecificName": "Swap",
											"readType": 0
										}`,
		expected: fmt.Errorf("invalid abi"),
	})

	testCases = append(testCases, testCase{
		name:                    "invalid read type",
		abiInput:                `{"constant":true,"inputs":[],"name":"someName","outputs":[],"payable":false,"stateMutability":"view","type":"function"}`,
		chainReadingDefinitions: `"TestMethod":{"readType": 59}`,
		expected:                fmt.Errorf("invalid chain reading definition read type: 59"),
	})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg, err := chainReaderTestHelper{}.makeChainReaderConfigFromStrings(tc.abiInput, tc.chainReadingDefinitions)
			assert.NoError(t, err)
			if tc.expected == nil {
				assert.NoError(t, validateChainReaderConfig(cfg))
			} else {
				assert.ErrorContains(t, validateChainReaderConfig(cfg), tc.expected.Error())
			}
		})
	}
}