// Code generated by mockery v2.46.0. DO NOT EDIT.

package date

import (
	time "time"

	mock "github.com/stretchr/testify/mock"
)

// MockHandler is an autogenerated mock type for the Handler type
type MockHandler struct {
	mock.Mock
}

type MockHandler_Expecter struct {
	mock *mock.Mock
}

func (_m *MockHandler) EXPECT() *MockHandler_Expecter {
	return &MockHandler_Expecter{mock: &_m.Mock}
}

// DateNow provides a mock function with given fields:
func (_m *MockHandler) DateNow() time.Time {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for DateNow")
	}

	var r0 time.Time
	if rf, ok := ret.Get(0).(func() time.Time); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Time)
	}

	return r0
}

// MockHandler_DateNow_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DateNow'
type MockHandler_DateNow_Call struct {
	*mock.Call
}

// DateNow is a helper method to define mock.On call
func (_e *MockHandler_Expecter) DateNow() *MockHandler_DateNow_Call {
	return &MockHandler_DateNow_Call{Call: _e.mock.On("DateNow")}
}

func (_c *MockHandler_DateNow_Call) Run(run func()) *MockHandler_DateNow_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockHandler_DateNow_Call) Return(_a0 time.Time) *MockHandler_DateNow_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockHandler_DateNow_Call) RunAndReturn(run func() time.Time) *MockHandler_DateNow_Call {
	_c.Call.Return(run)
	return _c
}

// DateParse provides a mock function with given fields: dateInSt
func (_m *MockHandler) DateParse(dateInSt string) (time.Time, error) {
	ret := _m.Called(dateInSt)

	if len(ret) == 0 {
		panic("no return value specified for DateParse")
	}

	var r0 time.Time
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (time.Time, error)); ok {
		return rf(dateInSt)
	}
	if rf, ok := ret.Get(0).(func(string) time.Time); ok {
		r0 = rf(dateInSt)
	} else {
		r0 = ret.Get(0).(time.Time)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(dateInSt)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockHandler_DateParse_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DateParse'
type MockHandler_DateParse_Call struct {
	*mock.Call
}

// DateParse is a helper method to define mock.On call
//   - dateInSt string
func (_e *MockHandler_Expecter) DateParse(dateInSt interface{}) *MockHandler_DateParse_Call {
	return &MockHandler_DateParse_Call{Call: _e.mock.On("DateParse", dateInSt)}
}

func (_c *MockHandler_DateParse_Call) Run(run func(dateInSt string)) *MockHandler_DateParse_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockHandler_DateParse_Call) Return(_a0 time.Time, _a1 error) *MockHandler_DateParse_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockHandler_DateParse_Call) RunAndReturn(run func(string) (time.Time, error)) *MockHandler_DateParse_Call {
	_c.Call.Return(run)
	return _c
}

// DateToString provides a mock function with given fields: _a0
func (_m *MockHandler) DateToString(_a0 time.Time) string {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for DateToString")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func(time.Time) string); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockHandler_DateToString_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DateToString'
type MockHandler_DateToString_Call struct {
	*mock.Call
}

// DateToString is a helper method to define mock.On call
//   - _a0 time.Time
func (_e *MockHandler_Expecter) DateToString(_a0 interface{}) *MockHandler_DateToString_Call {
	return &MockHandler_DateToString_Call{Call: _e.mock.On("DateToString", _a0)}
}

func (_c *MockHandler_DateToString_Call) Run(run func(_a0 time.Time)) *MockHandler_DateToString_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(time.Time))
	})
	return _c
}

func (_c *MockHandler_DateToString_Call) Return(_a0 string) *MockHandler_DateToString_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockHandler_DateToString_Call) RunAndReturn(run func(time.Time) string) *MockHandler_DateToString_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockHandler creates a new instance of MockHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockHandler {
	mock := &MockHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}