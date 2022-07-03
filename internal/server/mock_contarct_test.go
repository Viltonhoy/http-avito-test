// Code generated by MockGen. DO NOT EDIT.
// Source: contract.go

// Package server is a generated GoMock package.
package server

import (
	context "context"
	storage "http-avito-test/internal/storage"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	decimal "github.com/shopspring/decimal"
	zap "go.uber.org/zap"
)

// MockStorager is a mock of Storager interface.
type MockStorager struct {
	ctrl     *gomock.Controller
	recorder *MockStoragerMockRecorder
}

// MockStoragerMockRecorder is the mock recorder for MockStorager.
type MockStoragerMockRecorder struct {
	mock *MockStorager
}

// NewMockStorager creates a new mock instance.
func NewMockStorager(ctrl *gomock.Controller) *MockStorager {
	mock := &MockStorager{ctrl: ctrl}
	mock.recorder = &MockStoragerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorager) EXPECT() *MockStoragerMockRecorder {
	return m.recorder
}

// Deposit mocks base method.
func (m *MockStorager) Deposit(arg0 context.Context, arg1 int64, arg2 decimal.Decimal) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Deposit", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Deposit indicates an expected call of Deposit.
func (mr *MockStoragerMockRecorder) Deposit(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Deposit", reflect.TypeOf((*MockStorager)(nil).Deposit), arg0, arg1, arg2)
}

// ReadUserByID mocks base method.
func (m *MockStorager) ReadUserByID(arg0 context.Context, arg1 int64) (storage.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadUserByID", arg0, arg1)
	ret0, _ := ret[0].(storage.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadUserByID indicates an expected call of ReadUserByID.
func (mr *MockStoragerMockRecorder) ReadUserByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadUserByID", reflect.TypeOf((*MockStorager)(nil).ReadUserByID), arg0, arg1)
}

// ReadUserHistoryList mocks base method.
func (m *MockStorager) ReadUserHistoryList(ctx context.Context, user_id int64, order storage.OrdBy, limit, offset int64) ([]storage.ReadUserHistoryResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadUserHistoryList", ctx, user_id, order, limit, offset)
	ret0, _ := ret[0].([]storage.ReadUserHistoryResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadUserHistoryList indicates an expected call of ReadUserHistoryList.
func (mr *MockStoragerMockRecorder) ReadUserHistoryList(ctx, user_id, order, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadUserHistoryList", reflect.TypeOf((*MockStorager)(nil).ReadUserHistoryList), ctx, user_id, order, limit, offset)
}

// Transfer mocks base method.
func (m *MockStorager) Transfer(ctx context.Context, user_id1, user_id2 int64, amount decimal.Decimal, description *string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Transfer", ctx, user_id1, user_id2, amount, description)
	ret0, _ := ret[0].(error)
	return ret0
}

// Transfer indicates an expected call of Transfer.
func (mr *MockStoragerMockRecorder) Transfer(ctx, user_id1, user_id2, amount, description interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Transfer", reflect.TypeOf((*MockStorager)(nil).Transfer), ctx, user_id1, user_id2, amount, description)
}

// Withdrawal mocks base method.
func (m *MockStorager) Withdrawal(arg0 context.Context, arg1 int64, arg2 decimal.Decimal, arg3 *string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Withdrawal", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// Withdrawal indicates an expected call of Withdrawal.
func (mr *MockStoragerMockRecorder) Withdrawal(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Withdrawal", reflect.TypeOf((*MockStorager)(nil).Withdrawal), arg0, arg1, arg2, arg3)
}

// MockExchanger is a mock of Exchanger interface.
type MockExchanger struct {
	ctrl     *gomock.Controller
	recorder *MockExchangerMockRecorder
}

// MockExchangerMockRecorder is the mock recorder for MockExchanger.
type MockExchangerMockRecorder struct {
	mock *MockExchanger
}

// NewMockExchanger creates a new mock instance.
func NewMockExchanger(ctrl *gomock.Controller) *MockExchanger {
	mock := &MockExchanger{ctrl: ctrl}
	mock.recorder = &MockExchangerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockExchanger) EXPECT() *MockExchangerMockRecorder {
	return m.recorder
}

// ExchangeRates mocks base method.
func (m *MockExchanger) ExchangeRates(logger *zap.Logger, value decimal.Decimal, currency string) (decimal.Decimal, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExchangeRates", logger, value, currency)
	ret0, _ := ret[0].(decimal.Decimal)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExchangeRates indicates an expected call of ExchangeRates.
func (mr *MockExchangerMockRecorder) ExchangeRates(logger, value, currency interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExchangeRates", reflect.TypeOf((*MockExchanger)(nil).ExchangeRates), logger, value, currency)
}
