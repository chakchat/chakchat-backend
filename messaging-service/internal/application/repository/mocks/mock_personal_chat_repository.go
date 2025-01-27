// Code generated by mockery; DO NOT EDIT.
// github.com/vektra/mockery

package mocks

import (
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	mock "github.com/stretchr/testify/mock"
)

// NewMockPersonalChatRepository creates a new instance of MockPersonalChatRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockPersonalChatRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockPersonalChatRepository {
	mock := &MockPersonalChatRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// MockPersonalChatRepository is an autogenerated mock type for the PersonalChatRepository type
type MockPersonalChatRepository struct {
	mock.Mock
}

type MockPersonalChatRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockPersonalChatRepository) EXPECT() *MockPersonalChatRepository_Expecter {
	return &MockPersonalChatRepository_Expecter{mock: &_m.Mock}
}

// Create provides a mock function for the type MockPersonalChatRepository
func (_mock *MockPersonalChatRepository) Create(chat *domain.PersonalChat) (*domain.PersonalChat, error) {
	ret := _mock.Called(chat)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *domain.PersonalChat
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(*domain.PersonalChat) (*domain.PersonalChat, error)); ok {
		return returnFunc(chat)
	}
	if returnFunc, ok := ret.Get(0).(func(*domain.PersonalChat) *domain.PersonalChat); ok {
		r0 = returnFunc(chat)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.PersonalChat)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(*domain.PersonalChat) error); ok {
		r1 = returnFunc(chat)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockPersonalChatRepository_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type MockPersonalChatRepository_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - chat
func (_e *MockPersonalChatRepository_Expecter) Create(chat interface{}) *MockPersonalChatRepository_Create_Call {
	return &MockPersonalChatRepository_Create_Call{Call: _e.mock.On("Create", chat)}
}

func (_c *MockPersonalChatRepository_Create_Call) Run(run func(chat *domain.PersonalChat)) *MockPersonalChatRepository_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*domain.PersonalChat))
	})
	return _c
}

func (_c *MockPersonalChatRepository_Create_Call) Return(personalChat *domain.PersonalChat, err error) *MockPersonalChatRepository_Create_Call {
	_c.Call.Return(personalChat, err)
	return _c
}

func (_c *MockPersonalChatRepository_Create_Call) RunAndReturn(run func(chat *domain.PersonalChat) (*domain.PersonalChat, error)) *MockPersonalChatRepository_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function for the type MockPersonalChatRepository
func (_mock *MockPersonalChatRepository) Delete(d domain.ChatID) error {
	ret := _mock.Called(d)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(domain.ChatID) error); ok {
		r0 = returnFunc(d)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// MockPersonalChatRepository_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type MockPersonalChatRepository_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - d
func (_e *MockPersonalChatRepository_Expecter) Delete(d interface{}) *MockPersonalChatRepository_Delete_Call {
	return &MockPersonalChatRepository_Delete_Call{Call: _e.mock.On("Delete", d)}
}

func (_c *MockPersonalChatRepository_Delete_Call) Run(run func(d domain.ChatID)) *MockPersonalChatRepository_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(domain.ChatID))
	})
	return _c
}

func (_c *MockPersonalChatRepository_Delete_Call) Return(err error) *MockPersonalChatRepository_Delete_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockPersonalChatRepository_Delete_Call) RunAndReturn(run func(d domain.ChatID) error) *MockPersonalChatRepository_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// FindById provides a mock function for the type MockPersonalChatRepository
func (_mock *MockPersonalChatRepository) FindById(chatId domain.ChatID) (*domain.PersonalChat, error) {
	ret := _mock.Called(chatId)

	if len(ret) == 0 {
		panic("no return value specified for FindById")
	}

	var r0 *domain.PersonalChat
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(domain.ChatID) (*domain.PersonalChat, error)); ok {
		return returnFunc(chatId)
	}
	if returnFunc, ok := ret.Get(0).(func(domain.ChatID) *domain.PersonalChat); ok {
		r0 = returnFunc(chatId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.PersonalChat)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(domain.ChatID) error); ok {
		r1 = returnFunc(chatId)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockPersonalChatRepository_FindById_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindById'
type MockPersonalChatRepository_FindById_Call struct {
	*mock.Call
}

// FindById is a helper method to define mock.On call
//   - chatId
func (_e *MockPersonalChatRepository_Expecter) FindById(chatId interface{}) *MockPersonalChatRepository_FindById_Call {
	return &MockPersonalChatRepository_FindById_Call{Call: _e.mock.On("FindById", chatId)}
}

func (_c *MockPersonalChatRepository_FindById_Call) Run(run func(chatId domain.ChatID)) *MockPersonalChatRepository_FindById_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(domain.ChatID))
	})
	return _c
}

func (_c *MockPersonalChatRepository_FindById_Call) Return(personalChat *domain.PersonalChat, err error) *MockPersonalChatRepository_FindById_Call {
	_c.Call.Return(personalChat, err)
	return _c
}

func (_c *MockPersonalChatRepository_FindById_Call) RunAndReturn(run func(chatId domain.ChatID) (*domain.PersonalChat, error)) *MockPersonalChatRepository_FindById_Call {
	_c.Call.Return(run)
	return _c
}

// FindByMembers provides a mock function for the type MockPersonalChatRepository
func (_mock *MockPersonalChatRepository) FindByMembers(members [2]domain.UserID) (*domain.PersonalChat, error) {
	ret := _mock.Called(members)

	if len(ret) == 0 {
		panic("no return value specified for FindByMembers")
	}

	var r0 *domain.PersonalChat
	var r1 error
	if returnFunc, ok := ret.Get(0).(func([2]domain.UserID) (*domain.PersonalChat, error)); ok {
		return returnFunc(members)
	}
	if returnFunc, ok := ret.Get(0).(func([2]domain.UserID) *domain.PersonalChat); ok {
		r0 = returnFunc(members)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.PersonalChat)
		}
	}
	if returnFunc, ok := ret.Get(1).(func([2]domain.UserID) error); ok {
		r1 = returnFunc(members)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockPersonalChatRepository_FindByMembers_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindByMembers'
type MockPersonalChatRepository_FindByMembers_Call struct {
	*mock.Call
}

// FindByMembers is a helper method to define mock.On call
//   - members
func (_e *MockPersonalChatRepository_Expecter) FindByMembers(members interface{}) *MockPersonalChatRepository_FindByMembers_Call {
	return &MockPersonalChatRepository_FindByMembers_Call{Call: _e.mock.On("FindByMembers", members)}
}

func (_c *MockPersonalChatRepository_FindByMembers_Call) Run(run func(members [2]domain.UserID)) *MockPersonalChatRepository_FindByMembers_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([2]domain.UserID))
	})
	return _c
}

func (_c *MockPersonalChatRepository_FindByMembers_Call) Return(personalChat *domain.PersonalChat, err error) *MockPersonalChatRepository_FindByMembers_Call {
	_c.Call.Return(personalChat, err)
	return _c
}

func (_c *MockPersonalChatRepository_FindByMembers_Call) RunAndReturn(run func(members [2]domain.UserID) (*domain.PersonalChat, error)) *MockPersonalChatRepository_FindByMembers_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function for the type MockPersonalChatRepository
func (_mock *MockPersonalChatRepository) Update(personalChat *domain.PersonalChat) (*domain.PersonalChat, error) {
	ret := _mock.Called(personalChat)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 *domain.PersonalChat
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(*domain.PersonalChat) (*domain.PersonalChat, error)); ok {
		return returnFunc(personalChat)
	}
	if returnFunc, ok := ret.Get(0).(func(*domain.PersonalChat) *domain.PersonalChat); ok {
		r0 = returnFunc(personalChat)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.PersonalChat)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(*domain.PersonalChat) error); ok {
		r1 = returnFunc(personalChat)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockPersonalChatRepository_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type MockPersonalChatRepository_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - personalChat
func (_e *MockPersonalChatRepository_Expecter) Update(personalChat interface{}) *MockPersonalChatRepository_Update_Call {
	return &MockPersonalChatRepository_Update_Call{Call: _e.mock.On("Update", personalChat)}
}

func (_c *MockPersonalChatRepository_Update_Call) Run(run func(personalChat *domain.PersonalChat)) *MockPersonalChatRepository_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*domain.PersonalChat))
	})
	return _c
}

func (_c *MockPersonalChatRepository_Update_Call) Return(personalChat1 *domain.PersonalChat, err error) *MockPersonalChatRepository_Update_Call {
	_c.Call.Return(personalChat1, err)
	return _c
}

func (_c *MockPersonalChatRepository_Update_Call) RunAndReturn(run func(personalChat *domain.PersonalChat) (*domain.PersonalChat, error)) *MockPersonalChatRepository_Update_Call {
	_c.Call.Return(run)
	return _c
}
