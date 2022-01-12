package terratxm

import (
	"regexp"
	"testing"
	"time"

	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	tcmocks "github.com/smartcontractkit/chainlink-terra/pkg/terra/client/mocks"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	wasmtypes "github.com/terra-money/core/x/wasm/types"
)

func TestErrMatch(t *testing.T) {
	re, err := regexp.Compile(`^.*failed to execute message; message index: (?P<Index>\d{1}):.*$`)
	require.NoError(t, err)
	errStr := "rpc error: code = InvalidArgument desc = failed to execute message; message index: 0: Error parsing into type my_first_contract::msg::ExecuteMsg: unknown variant `blah`, expected `increment` or `reset`: execute wasm contract failed: invalid request"
	m := re.FindStringSubmatch(errStr)
	assert.Equal(t, m[1], "0")
	t.Log(len(m))
}

func TestTxm(t *testing.T) {
	t.Skip() // TODO: reenable
	db := pgtest.NewSqlxDB(t)
	lggr := logger.TestLogger(t)
	ks := keystore.New(db, utils.FastScryptParams, lggr, pgtest.NewPGCfg(true))
	require.NoError(t, ks.Unlock("blah"))
	k1, err := ks.Terra().Create()
	require.NoError(t, err)
	k2, err := ks.Terra().Create()
	require.NoError(t, err)

	t.Run("single msg", func(t *testing.T) {
		tc := new(tcmocks.ReaderWriter)
		tc.On("Account", mock.Anything).Return(uint64(0), uint64(0), nil)
		tc.On("GasPrice").Return(cosmostypes.NewDecCoinFromDec("uluna", cosmostypes.MustNewDecFromStr("0.01")))
		tc.On("SignAndBroadcast", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&txtypes.BroadcastTxResponse{
			TxResponse: &cosmostypes.TxResponse{TxHash: "0x123"},
		}, nil)
		tc.On("TxsEvents", mock.Anything).Return(&txtypes.GetTxsEventResponse{
			Txs:         []*txtypes.Tx{&txtypes.Tx{}},
			TxResponses: []*cosmostypes.TxResponse{{TxHash: "0x123"}},
		}, nil)

		txm := NewTxm(db, tc, ks.Terra(), lggr, pgtest.NewPGCfg(true), nil, time.Second)

		// Enqueue a single msg, then send it in a batch
		contract, err := cosmostypes.AccAddressFromBech32("terra1pp76d50yv2ldaahsdxdv8mmzqfjr2ax97gmue8")
		require.NoError(t, err)
		sender, err := cosmostypes.AccAddressFromBech32(k1.PublicKeyStr())
		require.NoError(t, err)
		msg1 := wasmtypes.NewMsgExecuteContract(sender, contract, []byte(`{"transmit":{"report_context":"","signatures":[""],"report":""}}`), cosmostypes.Coins{})
		d, err := msg1.Marshal()
		require.NoError(t, err)
		id1, err := txm.Enqueue(contract.String(), d)
		require.NoError(t, err)
		txm.sendMsgBatch()

		// Should be in completed state
		completed, err := txm.orm.SelectMsgsWithIDs([]int64{id1})
		require.NoError(t, err)
		require.Equal(t, 1, len(completed))
		assert.Equal(t, completed[0].State, Confirmed)
		tc.AssertExpectations(t)
	})

	t.Run("two msgs different accounts", func(t *testing.T) {
		tc := new(tcmocks.ReaderWriter)
		tc.On("Account", mock.Anything).Return(uint64(0), uint64(0), nil)
		tc.On("GasPrice").Return(cosmostypes.NewDecCoinFromDec("uluna", cosmostypes.MustNewDecFromStr("0.01")))
		tc.On("SignAndBroadcast", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&txtypes.BroadcastTxResponse{
			TxResponse: &cosmostypes.TxResponse{TxHash: "0x123"},
		}, nil)
		tc.On("TxsEvents", mock.Anything).Return(&txtypes.GetTxsEventResponse{
			Txs:         []*txtypes.Tx{&txtypes.Tx{}},
			TxResponses: []*cosmostypes.TxResponse{{TxHash: "0x123"}},
		}, nil)

		txm := NewTxm(db, tc, ks.Terra(), lggr, pgtest.NewPGCfg(true), nil, time.Second)

		contract, err := cosmostypes.AccAddressFromBech32("terra1pp76d50yv2ldaahsdxdv8mmzqfjr2ax97gmue8")
		require.NoError(t, err)
		sender1, err := cosmostypes.AccAddressFromBech32(k1.PublicKeyStr())
		require.NoError(t, err)
		msg1 := wasmtypes.NewMsgExecuteContract(sender1, contract, []byte(`{"transmit":{"report_context":"","signatures":[""],"report":""}}`), cosmostypes.Coins{})
		d, err := msg1.Marshal()
		require.NoError(t, err)

		sender2, err := cosmostypes.AccAddressFromBech32(k2.PublicKeyStr())
		require.NoError(t, err)
		msg2 := wasmtypes.NewMsgExecuteContract(sender2, contract, []byte(`{"transmit":{"report_context":"","signatures":[""],"report":""}}`), cosmostypes.Coins{})
		d2, err := msg2.Marshal()
		require.NoError(t, err)

		id1, err := txm.Enqueue(contract.String(), d)
		require.NoError(t, err)
		id2, err := txm.Enqueue(contract.String(), d2)
		require.NoError(t, err)
		txm.sendMsgBatch()

		// Should be in completed state
		completed, err := txm.orm.SelectMsgsWithIDs([]int64{id1, id2})
		require.NoError(t, err)
		require.Equal(t, 2, len(completed))
		assert.Equal(t, completed[0].State, Confirmed)
		assert.Equal(t, completed[1].State, Confirmed)
		tc.AssertExpectations(t)
	})
}
