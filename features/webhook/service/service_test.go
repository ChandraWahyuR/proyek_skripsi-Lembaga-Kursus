package service

import (
	"errors"
	"skripsi/features/transaksi/data"
	"skripsi/features/webhook"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockNotificationData struct {
	mock.Mock
}

func (m *MockNotificationData) HandleNotification(notification webhook.PaymentNotification, transaction data.Transaksi) error {
	args := m.Called(notification, transaction)
	return args.Error(0)
}

func TestWebhookService_HandleNotification(t *testing.T) {
	mockData := new(MockNotificationData)
	service := New(mockData)

	t.Run("success - transaction captured and accepted", func(t *testing.T) {
		notification := webhook.PaymentNotification{
			OrderID:           "12345",
			TransactionStatus: "capture",
			FraudStatus:       "accept",
		}
		expectedTransaction := data.Transaksi{
			ID:     "12345",
			Status: "Success",
		}

		mockData.On("HandleNotification", notification, expectedTransaction).Return(nil).Once()

		err := service.HandleNotification(notification)

		assert.Nil(t, err)
		mockData.AssertExpectations(t)
	})

	t.Run("success - transaction settled", func(t *testing.T) {
		notification := webhook.PaymentNotification{
			OrderID:           "67890",
			TransactionStatus: "settlement",
		}
		expectedTransaction := data.Transaksi{
			ID:     "67890",
			Status: "Success",
		}

		mockData.On("HandleNotification", notification, expectedTransaction).Return(nil).Once()

		err := service.HandleNotification(notification)

		assert.Nil(t, err)
		mockData.AssertExpectations(t)
	})

	t.Run("fail - transaction denied", func(t *testing.T) {
		notification := webhook.PaymentNotification{
			OrderID:           "54321",
			TransactionStatus: "deny",
		}
		expectedTransaction := data.Transaksi{
			ID:     "54321",
			Status: "Failed",
		}

		mockData.On("HandleNotification", notification, expectedTransaction).Return(nil).Once()

		err := service.HandleNotification(notification)

		assert.Nil(t, err)
		mockData.AssertExpectations(t)
	})

	t.Run("fail - data layer error", func(t *testing.T) {
		notification := webhook.PaymentNotification{
			OrderID:           "11223",
			TransactionStatus: "settlement",
		}
		expectedTransaction := data.Transaksi{
			ID:     "11223",
			Status: "Success",
		}
		mockData.On("HandleNotification", notification, expectedTransaction).Return(errors.New("data layer error")).Once()

		err := service.HandleNotification(notification)

		assert.NotNil(t, err)
		assert.Equal(t, "data layer error", err.Error())
		mockData.AssertExpectations(t)
	})
}
