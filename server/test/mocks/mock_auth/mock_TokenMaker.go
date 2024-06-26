// Code generated by mockery v2.42.3. DO NOT EDIT.

package mock_auth

import (
	time "time"

	model "github.com/dvid-messanger/internal/core/domain/model"
	mock "github.com/stretchr/testify/mock"
)

// MockTokenMaker is an autogenerated mock type for the TokenMaker type
type MockTokenMaker struct {
	mock.Mock
}

type MockTokenMaker_Expecter struct {
	mock *mock.Mock
}

func (_m *MockTokenMaker) EXPECT() *MockTokenMaker_Expecter {
	return &MockTokenMaker_Expecter{mock: &_m.Mock}
}

// MakeToken provides a mock function with given fields: user, duration
func (_m *MockTokenMaker) MakeToken(user model.UserCredentials, duration time.Duration) (string, error) {
	ret := _m.Called(user, duration)

	if len(ret) == 0 {
		panic("no return value specified for MakeToken")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(model.UserCredentials, time.Duration) (string, error)); ok {
		return rf(user, duration)
	}
	if rf, ok := ret.Get(0).(func(model.UserCredentials, time.Duration) string); ok {
		r0 = rf(user, duration)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(model.UserCredentials, time.Duration) error); ok {
		r1 = rf(user, duration)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockTokenMaker_MakeToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'MakeToken'
type MockTokenMaker_MakeToken_Call struct {
	*mock.Call
}

// MakeToken is a helper method to define mock.On call
//   - user model.UserCredentials
//   - duration time.Duration
func (_e *MockTokenMaker_Expecter) MakeToken(user interface{}, duration interface{}) *MockTokenMaker_MakeToken_Call {
	return &MockTokenMaker_MakeToken_Call{Call: _e.mock.On("MakeToken", user, duration)}
}

func (_c *MockTokenMaker_MakeToken_Call) Run(run func(user model.UserCredentials, duration time.Duration)) *MockTokenMaker_MakeToken_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(model.UserCredentials), args[1].(time.Duration))
	})
	return _c
}

func (_c *MockTokenMaker_MakeToken_Call) Return(_a0 string, _a1 error) *MockTokenMaker_MakeToken_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockTokenMaker_MakeToken_Call) RunAndReturn(run func(model.UserCredentials, time.Duration) (string, error)) *MockTokenMaker_MakeToken_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockTokenMaker creates a new instance of MockTokenMaker. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockTokenMaker(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockTokenMaker {
	mock := &MockTokenMaker{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
