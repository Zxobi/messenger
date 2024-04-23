// Code generated by mockery v2.42.3. DO NOT EDIT.

package mock_chat

import (
	context "context"

	model "github.com/dvid-messanger/internal/domain/model"
	mock "github.com/stretchr/testify/mock"
)

// MockChatProvider is an autogenerated mock type for the ChatProvider type
type MockChatProvider struct {
	mock.Mock
}

type MockChatProvider_Expecter struct {
	mock *mock.Mock
}

func (_m *MockChatProvider) EXPECT() *MockChatProvider_Expecter {
	return &MockChatProvider_Expecter{mock: &_m.Mock}
}

// Chat provides a mock function with given fields: ctx, cid
func (_m *MockChatProvider) Chat(ctx context.Context, cid []byte) (model.Chat, error) {
	ret := _m.Called(ctx, cid)

	if len(ret) == 0 {
		panic("no return value specified for Chat")
	}

	var r0 model.Chat
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []byte) (model.Chat, error)); ok {
		return rf(ctx, cid)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []byte) model.Chat); ok {
		r0 = rf(ctx, cid)
	} else {
		r0 = ret.Get(0).(model.Chat)
	}

	if rf, ok := ret.Get(1).(func(context.Context, []byte) error); ok {
		r1 = rf(ctx, cid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockChatProvider_Chat_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Chat'
type MockChatProvider_Chat_Call struct {
	*mock.Call
}

// Chat is a helper method to define mock.On call
//   - ctx context.Context
//   - cid []byte
func (_e *MockChatProvider_Expecter) Chat(ctx interface{}, cid interface{}) *MockChatProvider_Chat_Call {
	return &MockChatProvider_Chat_Call{Call: _e.mock.On("Chat", ctx, cid)}
}

func (_c *MockChatProvider_Chat_Call) Run(run func(ctx context.Context, cid []byte)) *MockChatProvider_Chat_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]byte))
	})
	return _c
}

func (_c *MockChatProvider_Chat_Call) Return(_a0 model.Chat, _a1 error) *MockChatProvider_Chat_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockChatProvider_Chat_Call) RunAndReturn(run func(context.Context, []byte) (model.Chat, error)) *MockChatProvider_Chat_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockChatProvider creates a new instance of MockChatProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockChatProvider(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockChatProvider {
	mock := &MockChatProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
