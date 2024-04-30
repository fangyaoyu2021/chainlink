// Code generated by mockery v2.42.2. DO NOT EDIT.

package mocks

import (
	big "math/big"

	assets "github.com/smartcontractkit/chainlink-common/pkg/assets"

	common "github.com/ethereum/go-ethereum/common"

	commontypes "github.com/smartcontractkit/chainlink/v2/common/types"

	context "context"

	coretypes "github.com/ethereum/go-ethereum/core/types"

	ethereum "github.com/ethereum/go-ethereum"

	evmassets "github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"

	mock "github.com/stretchr/testify/mock"

	rpc "github.com/ethereum/go-ethereum/rpc"

	types "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

// RPCClient is an autogenerated mock type for the RPCClient type
type RPCClient struct {
	mock.Mock
}

// BalanceAt provides a mock function with given fields: ctx, accountAddress, blockNumber
func (_m *RPCClient) BalanceAt(ctx context.Context, accountAddress common.Address, blockNumber *big.Int) (*big.Int, error) {
	ret := _m.Called(ctx, accountAddress, blockNumber)

	if len(ret) == 0 {
		panic("no return value specified for BalanceAt")
	}

	var r0 *big.Int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, common.Address, *big.Int) (*big.Int, error)); ok {
		return rf(ctx, accountAddress, blockNumber)
	}
	if rf, ok := ret.Get(0).(func(context.Context, common.Address, *big.Int) *big.Int); ok {
		r0 = rf(ctx, accountAddress, blockNumber)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, common.Address, *big.Int) error); ok {
		r1 = rf(ctx, accountAddress, blockNumber)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BatchCallContext provides a mock function with given fields: ctx, b
func (_m *RPCClient) BatchCallContext(ctx context.Context, b []rpc.BatchElem) error {
	ret := _m.Called(ctx, b)

	if len(ret) == 0 {
		panic("no return value specified for BatchCallContext")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []rpc.BatchElem) error); ok {
		r0 = rf(ctx, b)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// BlockByHash provides a mock function with given fields: ctx, hash
func (_m *RPCClient) BlockByHash(ctx context.Context, hash common.Hash) (*types.Head, error) {
	ret := _m.Called(ctx, hash)

	if len(ret) == 0 {
		panic("no return value specified for BlockByHash")
	}

	var r0 *types.Head
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, common.Hash) (*types.Head, error)); ok {
		return rf(ctx, hash)
	}
	if rf, ok := ret.Get(0).(func(context.Context, common.Hash) *types.Head); ok {
		r0 = rf(ctx, hash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Head)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, common.Hash) error); ok {
		r1 = rf(ctx, hash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BlockByHashGeth provides a mock function with given fields: ctx, hash
func (_m *RPCClient) BlockByHashGeth(ctx context.Context, hash common.Hash) (*coretypes.Block, error) {
	ret := _m.Called(ctx, hash)

	if len(ret) == 0 {
		panic("no return value specified for BlockByHashGeth")
	}

	var r0 *coretypes.Block
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, common.Hash) (*coretypes.Block, error)); ok {
		return rf(ctx, hash)
	}
	if rf, ok := ret.Get(0).(func(context.Context, common.Hash) *coretypes.Block); ok {
		r0 = rf(ctx, hash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*coretypes.Block)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, common.Hash) error); ok {
		r1 = rf(ctx, hash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BlockByNumber provides a mock function with given fields: ctx, number
func (_m *RPCClient) BlockByNumber(ctx context.Context, number *big.Int) (*types.Head, error) {
	ret := _m.Called(ctx, number)

	if len(ret) == 0 {
		panic("no return value specified for BlockByNumber")
	}

	var r0 *types.Head
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *big.Int) (*types.Head, error)); ok {
		return rf(ctx, number)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *big.Int) *types.Head); ok {
		r0 = rf(ctx, number)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Head)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *big.Int) error); ok {
		r1 = rf(ctx, number)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BlockByNumberGeth provides a mock function with given fields: ctx, number
func (_m *RPCClient) BlockByNumberGeth(ctx context.Context, number *big.Int) (*coretypes.Block, error) {
	ret := _m.Called(ctx, number)

	if len(ret) == 0 {
		panic("no return value specified for BlockByNumberGeth")
	}

	var r0 *coretypes.Block
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *big.Int) (*coretypes.Block, error)); ok {
		return rf(ctx, number)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *big.Int) *coretypes.Block); ok {
		r0 = rf(ctx, number)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*coretypes.Block)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *big.Int) error); ok {
		r1 = rf(ctx, number)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CallContext provides a mock function with given fields: ctx, result, method, args
func (_m *RPCClient) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	var _ca []interface{}
	_ca = append(_ca, ctx, result, method)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for CallContext")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, string, ...interface{}) error); ok {
		r0 = rf(ctx, result, method, args...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CallContract provides a mock function with given fields: ctx, msg, blockNumber
func (_m *RPCClient) CallContract(ctx context.Context, msg interface{}, blockNumber *big.Int) ([]byte, error) {
	ret := _m.Called(ctx, msg, blockNumber)

	if len(ret) == 0 {
		panic("no return value specified for CallContract")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, *big.Int) ([]byte, error)); ok {
		return rf(ctx, msg, blockNumber)
	}
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, *big.Int) []byte); ok {
		r0 = rf(ctx, msg, blockNumber)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, interface{}, *big.Int) error); ok {
		r1 = rf(ctx, msg, blockNumber)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ChainID provides a mock function with given fields: ctx
func (_m *RPCClient) ChainID(ctx context.Context) (*big.Int, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for ChainID")
	}

	var r0 *big.Int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*big.Int, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *big.Int); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ClientVersion provides a mock function with given fields: _a0
func (_m *RPCClient) ClientVersion(_a0 context.Context) (string, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for ClientVersion")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (string, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(context.Context) string); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Close provides a mock function with given fields:
func (_m *RPCClient) Close() {
	_m.Called()
}

// CodeAt provides a mock function with given fields: ctx, account, blockNumber
func (_m *RPCClient) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	ret := _m.Called(ctx, account, blockNumber)

	if len(ret) == 0 {
		panic("no return value specified for CodeAt")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, common.Address, *big.Int) ([]byte, error)); ok {
		return rf(ctx, account, blockNumber)
	}
	if rf, ok := ret.Get(0).(func(context.Context, common.Address, *big.Int) []byte); ok {
		r0 = rf(ctx, account, blockNumber)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, common.Address, *big.Int) error); ok {
		r1 = rf(ctx, account, blockNumber)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Dial provides a mock function with given fields: ctx
func (_m *RPCClient) Dial(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Dial")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DialHTTP provides a mock function with given fields:
func (_m *RPCClient) DialHTTP() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for DialHTTP")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DisconnectAll provides a mock function with given fields:
func (_m *RPCClient) DisconnectAll() {
	_m.Called()
}

// EstimateGas provides a mock function with given fields: ctx, call
func (_m *RPCClient) EstimateGas(ctx context.Context, call interface{}) (uint64, error) {
	ret := _m.Called(ctx, call)

	if len(ret) == 0 {
		panic("no return value specified for EstimateGas")
	}

	var r0 uint64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) (uint64, error)); ok {
		return rf(ctx, call)
	}
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) uint64); ok {
		r0 = rf(ctx, call)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, interface{}) error); ok {
		r1 = rf(ctx, call)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FilterEvents provides a mock function with given fields: ctx, query
func (_m *RPCClient) FilterEvents(ctx context.Context, query ethereum.FilterQuery) ([]coretypes.Log, error) {
	ret := _m.Called(ctx, query)

	if len(ret) == 0 {
		panic("no return value specified for FilterEvents")
	}

	var r0 []coretypes.Log
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ethereum.FilterQuery) ([]coretypes.Log, error)); ok {
		return rf(ctx, query)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ethereum.FilterQuery) []coretypes.Log); ok {
		r0 = rf(ctx, query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]coretypes.Log)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ethereum.FilterQuery) error); ok {
		r1 = rf(ctx, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetInterceptedChainInfo provides a mock function with given fields:
func (_m *RPCClient) GetInterceptedChainInfo() (int64, int64) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetInterceptedChainInfo")
	}

	var r0 int64
	var r1 int64
	if rf, ok := ret.Get(0).(func() (int64, int64)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func() int64); ok {
		r1 = rf()
	} else {
		r1 = ret.Get(1).(int64)
	}

	return r0, r1
}

// HeaderByHash provides a mock function with given fields: ctx, h
func (_m *RPCClient) HeaderByHash(ctx context.Context, h common.Hash) (*coretypes.Header, error) {
	ret := _m.Called(ctx, h)

	if len(ret) == 0 {
		panic("no return value specified for HeaderByHash")
	}

	var r0 *coretypes.Header
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, common.Hash) (*coretypes.Header, error)); ok {
		return rf(ctx, h)
	}
	if rf, ok := ret.Get(0).(func(context.Context, common.Hash) *coretypes.Header); ok {
		r0 = rf(ctx, h)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*coretypes.Header)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, common.Hash) error); ok {
		r1 = rf(ctx, h)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// HeaderByNumber provides a mock function with given fields: ctx, n
func (_m *RPCClient) HeaderByNumber(ctx context.Context, n *big.Int) (*coretypes.Header, error) {
	ret := _m.Called(ctx, n)

	if len(ret) == 0 {
		panic("no return value specified for HeaderByNumber")
	}

	var r0 *coretypes.Header
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *big.Int) (*coretypes.Header, error)); ok {
		return rf(ctx, n)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *big.Int) *coretypes.Header); ok {
		r0 = rf(ctx, n)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*coretypes.Header)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *big.Int) error); ok {
		r1 = rf(ctx, n)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsSyncing provides a mock function with given fields: ctx
func (_m *RPCClient) IsSyncing(ctx context.Context) (bool, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for IsSyncing")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (bool, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) bool); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LINKBalance provides a mock function with given fields: ctx, accountAddress, linkAddress
func (_m *RPCClient) LINKBalance(ctx context.Context, accountAddress common.Address, linkAddress common.Address) (*assets.Link, error) {
	ret := _m.Called(ctx, accountAddress, linkAddress)

	if len(ret) == 0 {
		panic("no return value specified for LINKBalance")
	}

	var r0 *assets.Link
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, common.Address, common.Address) (*assets.Link, error)); ok {
		return rf(ctx, accountAddress, linkAddress)
	}
	if rf, ok := ret.Get(0).(func(context.Context, common.Address, common.Address) *assets.Link); ok {
		r0 = rf(ctx, accountAddress, linkAddress)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*assets.Link)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, common.Address, common.Address) error); ok {
		r1 = rf(ctx, accountAddress, linkAddress)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LatestBlockHeight provides a mock function with given fields: _a0
func (_m *RPCClient) LatestBlockHeight(_a0 context.Context) (*big.Int, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for LatestBlockHeight")
	}

	var r0 *big.Int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*big.Int, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *big.Int); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LatestFinalizedBlock provides a mock function with given fields: ctx
func (_m *RPCClient) LatestFinalizedBlock(ctx context.Context) (*types.Head, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for LatestFinalizedBlock")
	}

	var r0 *types.Head
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*types.Head, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *types.Head); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Head)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PendingCallContract provides a mock function with given fields: ctx, msg
func (_m *RPCClient) PendingCallContract(ctx context.Context, msg interface{}) ([]byte, error) {
	ret := _m.Called(ctx, msg)

	if len(ret) == 0 {
		panic("no return value specified for PendingCallContract")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) ([]byte, error)); ok {
		return rf(ctx, msg)
	}
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) []byte); ok {
		r0 = rf(ctx, msg)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, interface{}) error); ok {
		r1 = rf(ctx, msg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PendingCodeAt provides a mock function with given fields: ctx, account
func (_m *RPCClient) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	ret := _m.Called(ctx, account)

	if len(ret) == 0 {
		panic("no return value specified for PendingCodeAt")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, common.Address) ([]byte, error)); ok {
		return rf(ctx, account)
	}
	if rf, ok := ret.Get(0).(func(context.Context, common.Address) []byte); ok {
		r0 = rf(ctx, account)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, common.Address) error); ok {
		r1 = rf(ctx, account)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PendingSequenceAt provides a mock function with given fields: ctx, addr
func (_m *RPCClient) PendingSequenceAt(ctx context.Context, addr common.Address) (types.Nonce, error) {
	ret := _m.Called(ctx, addr)

	if len(ret) == 0 {
		panic("no return value specified for PendingSequenceAt")
	}

	var r0 types.Nonce
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, common.Address) (types.Nonce, error)); ok {
		return rf(ctx, addr)
	}
	if rf, ok := ret.Get(0).(func(context.Context, common.Address) types.Nonce); ok {
		r0 = rf(ctx, addr)
	} else {
		r0 = ret.Get(0).(types.Nonce)
	}

	if rf, ok := ret.Get(1).(func(context.Context, common.Address) error); ok {
		r1 = rf(ctx, addr)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SendEmptyTransaction provides a mock function with given fields: ctx, newTxAttempt, seq, gasLimit, fee, fromAddress
func (_m *RPCClient) SendEmptyTransaction(ctx context.Context, newTxAttempt func(types.Nonce, uint32, *evmassets.Wei, common.Address) (interface{}, error), seq types.Nonce, gasLimit uint32, fee *evmassets.Wei, fromAddress common.Address) (string, error) {
	ret := _m.Called(ctx, newTxAttempt, seq, gasLimit, fee, fromAddress)

	if len(ret) == 0 {
		panic("no return value specified for SendEmptyTransaction")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, func(types.Nonce, uint32, *evmassets.Wei, common.Address) (interface{}, error), types.Nonce, uint32, *evmassets.Wei, common.Address) (string, error)); ok {
		return rf(ctx, newTxAttempt, seq, gasLimit, fee, fromAddress)
	}
	if rf, ok := ret.Get(0).(func(context.Context, func(types.Nonce, uint32, *evmassets.Wei, common.Address) (interface{}, error), types.Nonce, uint32, *evmassets.Wei, common.Address) string); ok {
		r0 = rf(ctx, newTxAttempt, seq, gasLimit, fee, fromAddress)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, func(types.Nonce, uint32, *evmassets.Wei, common.Address) (interface{}, error), types.Nonce, uint32, *evmassets.Wei, common.Address) error); ok {
		r1 = rf(ctx, newTxAttempt, seq, gasLimit, fee, fromAddress)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SendTransaction provides a mock function with given fields: ctx, tx
func (_m *RPCClient) SendTransaction(ctx context.Context, tx *coretypes.Transaction) error {
	ret := _m.Called(ctx, tx)

	if len(ret) == 0 {
		panic("no return value specified for SendTransaction")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *coretypes.Transaction) error); ok {
		r0 = rf(ctx, tx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SequenceAt provides a mock function with given fields: ctx, accountAddress, blockNumber
func (_m *RPCClient) SequenceAt(ctx context.Context, accountAddress common.Address, blockNumber *big.Int) (types.Nonce, error) {
	ret := _m.Called(ctx, accountAddress, blockNumber)

	if len(ret) == 0 {
		panic("no return value specified for SequenceAt")
	}

	var r0 types.Nonce
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, common.Address, *big.Int) (types.Nonce, error)); ok {
		return rf(ctx, accountAddress, blockNumber)
	}
	if rf, ok := ret.Get(0).(func(context.Context, common.Address, *big.Int) types.Nonce); ok {
		r0 = rf(ctx, accountAddress, blockNumber)
	} else {
		r0 = ret.Get(0).(types.Nonce)
	}

	if rf, ok := ret.Get(1).(func(context.Context, common.Address, *big.Int) error); ok {
		r1 = rf(ctx, accountAddress, blockNumber)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetAliveLoopSub provides a mock function with given fields: _a0
func (_m *RPCClient) SetAliveLoopSub(_a0 commontypes.Subscription) {
	_m.Called(_a0)
}

// SimulateTransaction provides a mock function with given fields: ctx, tx
func (_m *RPCClient) SimulateTransaction(ctx context.Context, tx *coretypes.Transaction) error {
	ret := _m.Called(ctx, tx)

	if len(ret) == 0 {
		panic("no return value specified for SimulateTransaction")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *coretypes.Transaction) error); ok {
		r0 = rf(ctx, tx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SubscribeFilterLogs provides a mock function with given fields: ctx, q, ch
func (_m *RPCClient) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- coretypes.Log) (ethereum.Subscription, error) {
	ret := _m.Called(ctx, q, ch)

	if len(ret) == 0 {
		panic("no return value specified for SubscribeFilterLogs")
	}

	var r0 ethereum.Subscription
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ethereum.FilterQuery, chan<- coretypes.Log) (ethereum.Subscription, error)); ok {
		return rf(ctx, q, ch)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ethereum.FilterQuery, chan<- coretypes.Log) ethereum.Subscription); ok {
		r0 = rf(ctx, q, ch)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(ethereum.Subscription)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ethereum.FilterQuery, chan<- coretypes.Log) error); ok {
		r1 = rf(ctx, q, ch)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SubscribeNewHead provides a mock function with given fields: ctx, channel
func (_m *RPCClient) SubscribeNewHead(ctx context.Context, channel chan<- *types.Head) (commontypes.Subscription, error) {
	ret := _m.Called(ctx, channel)

	if len(ret) == 0 {
		panic("no return value specified for SubscribeNewHead")
	}

	var r0 commontypes.Subscription
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, chan<- *types.Head) (commontypes.Subscription, error)); ok {
		return rf(ctx, channel)
	}
	if rf, ok := ret.Get(0).(func(context.Context, chan<- *types.Head) commontypes.Subscription); ok {
		r0 = rf(ctx, channel)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(commontypes.Subscription)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, chan<- *types.Head) error); ok {
		r1 = rf(ctx, channel)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SubscribersCount provides a mock function with given fields:
func (_m *RPCClient) SubscribersCount() int32 {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for SubscribersCount")
	}

	var r0 int32
	if rf, ok := ret.Get(0).(func() int32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int32)
	}

	return r0
}

// SuggestGasPrice provides a mock function with given fields: ctx
func (_m *RPCClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for SuggestGasPrice")
	}

	var r0 *big.Int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*big.Int, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *big.Int); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SuggestGasTipCap provides a mock function with given fields: ctx
func (_m *RPCClient) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for SuggestGasTipCap")
	}

	var r0 *big.Int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*big.Int, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *big.Int); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TokenBalance provides a mock function with given fields: ctx, accountAddress, tokenAddress
func (_m *RPCClient) TokenBalance(ctx context.Context, accountAddress common.Address, tokenAddress common.Address) (*big.Int, error) {
	ret := _m.Called(ctx, accountAddress, tokenAddress)

	if len(ret) == 0 {
		panic("no return value specified for TokenBalance")
	}

	var r0 *big.Int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, common.Address, common.Address) (*big.Int, error)); ok {
		return rf(ctx, accountAddress, tokenAddress)
	}
	if rf, ok := ret.Get(0).(func(context.Context, common.Address, common.Address) *big.Int); ok {
		r0 = rf(ctx, accountAddress, tokenAddress)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, common.Address, common.Address) error); ok {
		r1 = rf(ctx, accountAddress, tokenAddress)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TransactionByHash provides a mock function with given fields: ctx, txHash
func (_m *RPCClient) TransactionByHash(ctx context.Context, txHash common.Hash) (*coretypes.Transaction, error) {
	ret := _m.Called(ctx, txHash)

	if len(ret) == 0 {
		panic("no return value specified for TransactionByHash")
	}

	var r0 *coretypes.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, common.Hash) (*coretypes.Transaction, error)); ok {
		return rf(ctx, txHash)
	}
	if rf, ok := ret.Get(0).(func(context.Context, common.Hash) *coretypes.Transaction); ok {
		r0 = rf(ctx, txHash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*coretypes.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, common.Hash) error); ok {
		r1 = rf(ctx, txHash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TransactionReceipt provides a mock function with given fields: ctx, txHash
func (_m *RPCClient) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	ret := _m.Called(ctx, txHash)

	if len(ret) == 0 {
		panic("no return value specified for TransactionReceipt")
	}

	var r0 *types.Receipt
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, common.Hash) (*types.Receipt, error)); ok {
		return rf(ctx, txHash)
	}
	if rf, ok := ret.Get(0).(func(context.Context, common.Hash) *types.Receipt); ok {
		r0 = rf(ctx, txHash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Receipt)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, common.Hash) error); ok {
		r1 = rf(ctx, txHash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TransactionReceiptGeth provides a mock function with given fields: ctx, txHash
func (_m *RPCClient) TransactionReceiptGeth(ctx context.Context, txHash common.Hash) (*coretypes.Receipt, error) {
	ret := _m.Called(ctx, txHash)

	if len(ret) == 0 {
		panic("no return value specified for TransactionReceiptGeth")
	}

	var r0 *coretypes.Receipt
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, common.Hash) (*coretypes.Receipt, error)); ok {
		return rf(ctx, txHash)
	}
	if rf, ok := ret.Get(0).(func(context.Context, common.Hash) *coretypes.Receipt); ok {
		r0 = rf(ctx, txHash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*coretypes.Receipt)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, common.Hash) error); ok {
		r1 = rf(ctx, txHash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UnsubscribeAllExceptAliveLoop provides a mock function with given fields:
func (_m *RPCClient) UnsubscribeAllExceptAliveLoop() {
	_m.Called()
}

// NewRPCClient creates a new instance of RPCClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRPCClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *RPCClient {
	mock := &RPCClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
