// Code generated by mockery v2.46.0. DO NOT EDIT.

package mocks

import (
	context "context"
	internal "library-api/internal"

	mock "github.com/stretchr/testify/mock"

	server "library-api/internal/platform/server"
)

// BookAdapter is an autogenerated mock type for the BookAdapter type
type BookAdapter struct {
	mock.Mock
}

type BookAdapter_Expecter struct {
	mock *mock.Mock
}

func (_m *BookAdapter) EXPECT() *BookAdapter_Expecter {
	return &BookAdapter_Expecter{mock: &_m.Mock}
}

// CreateBook provides a mock function with given fields: ctx, request
func (_m *BookAdapter) CreateBook(ctx context.Context, request server.CreateBookRequest) (internal.Book, error) {
	ret := _m.Called(ctx, request)

	if len(ret) == 0 {
		panic("no return value specified for CreateBook")
	}

	var r0 internal.Book
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, server.CreateBookRequest) (internal.Book, error)); ok {
		return rf(ctx, request)
	}
	if rf, ok := ret.Get(0).(func(context.Context, server.CreateBookRequest) internal.Book); ok {
		r0 = rf(ctx, request)
	} else {
		r0 = ret.Get(0).(internal.Book)
	}

	if rf, ok := ret.Get(1).(func(context.Context, server.CreateBookRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BookAdapter_CreateBook_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateBook'
type BookAdapter_CreateBook_Call struct {
	*mock.Call
}

// CreateBook is a helper method to define mock.On call
//   - ctx context.Context
//   - request server.CreateBookRequest
func (_e *BookAdapter_Expecter) CreateBook(ctx interface{}, request interface{}) *BookAdapter_CreateBook_Call {
	return &BookAdapter_CreateBook_Call{Call: _e.mock.On("CreateBook", ctx, request)}
}

func (_c *BookAdapter_CreateBook_Call) Run(run func(ctx context.Context, request server.CreateBookRequest)) *BookAdapter_CreateBook_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(server.CreateBookRequest))
	})
	return _c
}

func (_c *BookAdapter_CreateBook_Call) Return(_a0 internal.Book, _a1 error) *BookAdapter_CreateBook_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *BookAdapter_CreateBook_Call) RunAndReturn(run func(context.Context, server.CreateBookRequest) (internal.Book, error)) *BookAdapter_CreateBook_Call {
	_c.Call.Return(run)
	return _c
}

// GetBookByID provides a mock function with given fields: ctx, id
func (_m *BookAdapter) GetBookByID(ctx context.Context, id string) (internal.Book, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetBookByID")
	}

	var r0 internal.Book
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (internal.Book, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) internal.Book); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(internal.Book)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BookAdapter_GetBookByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetBookByID'
type BookAdapter_GetBookByID_Call struct {
	*mock.Call
}

// GetBookByID is a helper method to define mock.On call
//   - ctx context.Context
//   - id string
func (_e *BookAdapter_Expecter) GetBookByID(ctx interface{}, id interface{}) *BookAdapter_GetBookByID_Call {
	return &BookAdapter_GetBookByID_Call{Call: _e.mock.On("GetBookByID", ctx, id)}
}

func (_c *BookAdapter_GetBookByID_Call) Run(run func(ctx context.Context, id string)) *BookAdapter_GetBookByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *BookAdapter_GetBookByID_Call) Return(_a0 internal.Book, _a1 error) *BookAdapter_GetBookByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *BookAdapter_GetBookByID_Call) RunAndReturn(run func(context.Context, string) (internal.Book, error)) *BookAdapter_GetBookByID_Call {
	_c.Call.Return(run)
	return _c
}

// GetBooks provides a mock function with given fields: ctx
func (_m *BookAdapter) GetBooks(ctx context.Context) ([]internal.Book, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetBooks")
	}

	var r0 []internal.Book
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]internal.Book, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []internal.Book); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]internal.Book)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BookAdapter_GetBooks_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetBooks'
type BookAdapter_GetBooks_Call struct {
	*mock.Call
}

// GetBooks is a helper method to define mock.On call
//   - ctx context.Context
func (_e *BookAdapter_Expecter) GetBooks(ctx interface{}) *BookAdapter_GetBooks_Call {
	return &BookAdapter_GetBooks_Call{Call: _e.mock.On("GetBooks", ctx)}
}

func (_c *BookAdapter_GetBooks_Call) Run(run func(ctx context.Context)) *BookAdapter_GetBooks_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *BookAdapter_GetBooks_Call) Return(_a0 []internal.Book, _a1 error) *BookAdapter_GetBooks_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *BookAdapter_GetBooks_Call) RunAndReturn(run func(context.Context) ([]internal.Book, error)) *BookAdapter_GetBooks_Call {
	_c.Call.Return(run)
	return _c
}

// NewBookAdapter creates a new instance of BookAdapter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewBookAdapter(t interface {
	mock.TestingT
	Cleanup(func())
}) *BookAdapter {
	mock := &BookAdapter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
