// Code generated by mockery v2.10.1. DO NOT EDIT.

package mocks

import (
	big "math/big"

	assets "github.com/smartcontractkit/chainlink/core/assets"

	common "github.com/ethereum/go-ethereum/common"

	context "context"

	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"

	gas "github.com/smartcontractkit/chainlink/core/chains/evm/gas"

	mock "github.com/stretchr/testify/mock"

	pg "github.com/smartcontractkit/chainlink/core/services/pg"

	txmgr "github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"

	types "github.com/smartcontractkit/chainlink/common/txmgr/types"
)

// TxManager is an autogenerated mock type for the TxManager type
type TxManager struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *TxManager) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateEthTransaction provides a mock function with given fields: newTx, qopts
func (_m *TxManager) CreateEthTransaction(newTx txmgr.NewTx, qopts ...pg.QOpt) (txmgr.EthTx, error) {
	_va := make([]interface{}, len(qopts))
	for _i := range qopts {
		_va[_i] = qopts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, newTx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 txmgr.EthTx
	if rf, ok := ret.Get(0).(func(txmgr.NewTx, ...pg.QOpt) txmgr.EthTx); ok {
		r0 = rf(newTx, qopts...)
	} else {
		r0 = ret.Get(0).(txmgr.EthTx)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(txmgr.NewTx, ...pg.QOpt) error); ok {
		r1 = rf(newTx, qopts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetForwarderForEOA provides a mock function with given fields: eoa
func (_m *TxManager) GetForwarderForEOA(eoa common.Address) (common.Address, error) {
	ret := _m.Called(eoa)

	var r0 common.Address
	if rf, ok := ret.Get(0).(func(common.Address) common.Address); ok {
		r0 = rf(eoa)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.Address)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Address) error); ok {
		r1 = rf(eoa)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetGasEstimator provides a mock function with given fields:
func (_m *TxManager) GetGasEstimator() types.FeeEstimator[*evmtypes.Head, gas.EvmFee, *assets.Wei, common.Hash] {
	ret := _m.Called()

	var r0 types.FeeEstimator[*evmtypes.Head, gas.EvmFee, *assets.Wei, common.Hash]
	if rf, ok := ret.Get(0).(func() types.FeeEstimator[*evmtypes.Head, gas.EvmFee, *assets.Wei, common.Hash]); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.FeeEstimator[*evmtypes.Head, gas.EvmFee, *assets.Wei, common.Hash])
		}
	}

	return r0
}

// HealthReport provides a mock function with given fields:
func (_m *TxManager) HealthReport() map[string]error {
	ret := _m.Called()

	var r0 map[string]error
	if rf, ok := ret.Get(0).(func() map[string]error); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]error)
		}
	}

	return r0
}

// Healthy provides a mock function with given fields:
func (_m *TxManager) Healthy() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Name provides a mock function with given fields:
func (_m *TxManager) Name() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// OnNewLongestChain provides a mock function with given fields: ctx, head
<<<<<<< HEAD
func (_m *TxManager) OnNewLongestChain(ctx context.Context, head types.HeadView) {
=======
func (_m *TxManager) OnNewLongestChain(ctx context.Context, head *evmtypes.Head) {
>>>>>>> c90fc0ed3d2d8e1ad0b021e32ced366c040ab1f8
	_m.Called(ctx, head)
}

// Ready provides a mock function with given fields:
func (_m *TxManager) Ready() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RegisterResumeCallback provides a mock function with given fields: fn
func (_m *TxManager) RegisterResumeCallback(fn txmgr.ResumeCallback) {
	_m.Called(fn)
}

// Reset provides a mock function with given fields: f, addr, abandon
func (_m *TxManager) Reset(f func(), addr common.Address, abandon bool) error {
	ret := _m.Called(f, addr, abandon)

	var r0 error
	if rf, ok := ret.Get(0).(func(func(), common.Address, bool) error); ok {
		r0 = rf(f, addr, abandon)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SendEther provides a mock function with given fields: chainID, from, to, value, gasLimit
func (_m *TxManager) SendEther(chainID *big.Int, from common.Address, to common.Address, value assets.Eth, gasLimit uint32) (txmgr.EthTx, error) {
	ret := _m.Called(chainID, from, to, value, gasLimit)

	var r0 txmgr.EthTx
	if rf, ok := ret.Get(0).(func(*big.Int, common.Address, common.Address, assets.Eth, uint32) txmgr.EthTx); ok {
		r0 = rf(chainID, from, to, value, gasLimit)
	} else {
		r0 = ret.Get(0).(txmgr.EthTx)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*big.Int, common.Address, common.Address, assets.Eth, uint32) error); ok {
		r1 = rf(chainID, from, to, value, gasLimit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Start provides a mock function with given fields: _a0
func (_m *TxManager) Start(_a0 context.Context) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Trigger provides a mock function with given fields: addr
func (_m *TxManager) Trigger(addr common.Address) {
	_m.Called(addr)
}

type mockConstructorTestingTNewTxManager interface {
	mock.TestingT
	Cleanup(func())
}

// NewTxManager creates a new instance of TxManager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTxManager(t mockConstructorTestingTNewTxManager) *TxManager {
	mock := &TxManager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
