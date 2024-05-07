// Code generated by mockery v2.42.3. DO NOT EDIT.

package mock_auth

import (
	context "context"

	model "github.com/dvid-messanger/internal/core/domain/model"
	mock "github.com/stretchr/testify/mock"
)

// MockUserSaver is an autogenerated mock type for the UserSaver type
type MockUserSaver struct {
	mock.Mock
}

type MockUserSaver_Expecter struct {
	mock *mock.Mock
}

func (_m *MockUserSaver) EXPECT() *MockUserSaver_Expecter {
	return &MockUserSaver_Expecter{mock: &_m.Mock}
}

// Save provides a mock function with given fields: ctx, uid, email, passHash
func (_m *MockUserSaver) Save(ctx context.Context, uid []byte, email string, passHash []byte) (model.UserCredentials, error) {
	ret := _m.Called(ctx, uid, email, passHash)

	if len(ret) == 0 {
		panic("no return value specified for Save")
	}

	var r0 model.UserCredentials
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []byte, string, []byte) (model.UserCredentials, error)); ok {
		return rf(ctx, uid, email, passHash)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []byte, string, []byte) model.UserCredentials); ok {
		r0 = rf(ctx, uid, email, passHash)
	} else {
		r0 = ret.Get(0).(model.UserCredentials)
	}

	if rf, ok := ret.Get(1).(func(context.Context, []byte, string, []byte) error); ok {
		r1 = rf(ctx, uid, email, passHash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockUserSaver_Save_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Save'
type MockUserSaver_Save_Call struct {
	*mock.Call
}

// Save is a helper method to define mock.On call
//   - ctx context.Context
//   - uid []byte
//   - email string
//   - passHash []byte
func (_e *MockUserSaver_Expecter) Save(ctx interface{}, uid interface{}, email interface{}, passHash interface{}) *MockUserSaver_Save_Call {
	return &MockUserSaver_Save_Call{Call: _e.mock.On("Save", ctx, uid, email, passHash)}
}

func (_c *MockUserSaver_Save_Call) Run(run func(ctx context.Context, uid []byte, email string, passHash []byte)) *MockUserSaver_Save_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]byte), args[2].(string), args[3].([]byte))
	})
	return _c
}

func (_c *MockUserSaver_Save_Call) Return(_a0 model.UserCredentials, _a1 error) *MockUserSaver_Save_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockUserSaver_Save_Call) RunAndReturn(run func(context.Context, []byte, string, []byte) (model.UserCredentials, error)) *MockUserSaver_Save_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockUserSaver creates a new instance of MockUserSaver. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockUserSaver(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockUserSaver {
	mock := &MockUserSaver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}