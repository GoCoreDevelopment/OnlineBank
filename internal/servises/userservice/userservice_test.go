package userservice

import (
	"errors"
	"testing"

	"api/internal/models/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateUser(user user.UserRegistred) (int, error) {
	return 0, nil
}

func (m *MockRepository) GetUserPassword(email string) (string, error) {
	return "", nil
}

func (m *MockRepository) GetUserBalance(id int) (int, error) {
	return 0, nil
}

func (m *MockRepository) GetUserID(emain string) (int, error) {
	return 0, nil
}

func (m *MockRepository) TransferMoney(senderID, receiverID, amount int) error {
	args := m.Called(senderID, receiverID, amount)
	return args.Error(0)
}

func TestUserService_TransferBalance(t *testing.T) {
	tests := []struct {
		name string
		senderID int
		receiverID int
		amount int
		mockSetup func(m *MockRepository)
		expectedError error
	}{
		{
			name: "Succsess",
			senderID: 1,
			receiverID: 2,
			amount: 100,
			mockSetup: func(m *MockRepository) {
				m.On("TransferMoney", 1, 2, 100).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Amount zero",
			senderID: 1,
			receiverID: 2,
			amount: 0,
			mockSetup: func(m *MockRepository) {},
			expectedError: errors.New("the number cannot be zero"),
		},
		{
			name: "Check error sender not found",
			senderID: 1,
			receiverID: 2,
			amount: 100,
			mockSetup: func(m *MockRepository) {
				m.On("TransferMoney", 1, 2, 100).Return(errors.New("sender not found"))
			},
			expectedError: errors.New("sender not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockRepository{}
			tt.mockSetup(mockRepo)
			service := NewUserService(mockRepo)

			err := service.Transfer(tt.senderID, tt.receiverID, tt.amount)
			assert.Equal(t, tt.expectedError, err)

			mockRepo.AssertExpectations(t)
		})
	}
}