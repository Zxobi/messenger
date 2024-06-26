// Code generated by mockery v2.42.3. DO NOT EDIT.

package mock_user

import (
	context "context"

	model "github.com/dvid-messanger/internal/core/domain/model"
	mock "github.com/stretchr/testify/mock"
)

// MockUserProvider is an autogenerated mock type for the UserProvider type
type MockUserProvider struct {
	mock.Mock
}

type MockUserProvider_Expecter struct {
	mock *mock.Mock
}

func (_m *MockUserProvider) EXPECT() *MockUserProvider_Expecter {
	return &MockUserProvider_Expecter{mock: &_m.Mock}
}

// User provides a mock function with given fields: ctx, uid
func (_m *MockUserProvider) User(ctx context.Context, uid []byte) (model.User, error) {
	ret := _m.Called(ctx, uid)

	if len(ret) == 0 {
		panic("no return value specified for User")
	}

	var r0 model.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []byte) (model.User, error)); ok {
		return rf(ctx, uid)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []byte) model.User); ok {
		r0 = rf(ctx, uid)
	} else {
		r0 = ret.Get(0).(model.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, []byte) error); ok {
		r1 = rf(ctx, uid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockUserProvider_User_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'User'
type MockUserProvider_User_Call struct {
	*mock.Call
}

// User is a helper method to define mock.On call
//   - ctx context.Context
//   - uid []byte
func (_e *MockUserProvider_Expecter) User(ctx interface{}, uid interface{}) *MockUserProvider_User_Call {
	return &MockUserProvider_User_Call{Call: _e.mock.On("User", ctx, uid)}
}

func (_c *MockUserProvider_User_Call) Run(run func(ctx context.Context, uid []byte)) *MockUserProvider_User_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]byte))
	})
	return _c
}

func (_c *MockUserProvider_User_Call) Return(_a0 model.User, _a1 error) *MockUserProvider_User_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockUserProvider_User_Call) RunAndReturn(run func(context.Context, []byte) (model.User, error)) *MockUserProvider_User_Call {
	_c.Call.Return(run)
	return _c
}

// Users provides a mock function with given fields: ctx
func (_m *MockUserProvider) Users(ctx context.Context) ([]model.User, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Users")
	}

	var r0 []model.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]model.User, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []model.User); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockUserProvider_Users_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Users'
type MockUserProvider_Users_Call struct {
	*mock.Call
}

// Users is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockUserProvider_Expecter) Users(ctx interface{}) *MockUserProvider_Users_Call {
	return &MockUserProvider_Users_Call{Call: _e.mock.On("Users", ctx)}
}

func (_c *MockUserProvider_Users_Call) Run(run func(ctx context.Context)) *MockUserProvider_Users_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockUserProvider_Users_Call) Return(_a0 []model.User, _a1 error) *MockUserProvider_Users_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockUserProvider_Users_Call) RunAndReturn(run func(context.Context) ([]model.User, error)) *MockUserProvider_Users_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockUserProvider creates a new instance of MockUserProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockUserProvider(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockUserProvider {
	mock := &MockUserProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
