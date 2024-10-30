package service

import (
	"fmt"
	transaksiData "skripsi/features/transaksi/data"
	"skripsi/features/webhook"
)

type WebhookService struct {
	d webhook.MidtransNotificationData
}

func New(data webhook.MidtransNotificationData) webhook.MidtransNotificationService {
	return &WebhookService{
		d: data,
	}
}

func (s *WebhookService) HandleNotification(notification webhook.PaymentNotification) error {
	transactionStatus := notification.TransactionStatus
	fraudstatus := notification.FraudStatus

	transactionData := transaksiData.Transaksi{
		ID: notification.OrderID,
	}

	switch transactionStatus {
	case "capture":
		if fraudstatus == "accept" {
			transactionData.Status = "Success"
		}
	case "settlement":
		transactionData.Status = "Success"
	case "cancel", "deny", "expire":
		transactionData.Status = "Failed"
	case "pending":
		transactionData.Status = "Pending"
	}

	fmt.Printf("Processing notification: %+v\n", notification)
	fmt.Printf("Updating transaction data: %+v\n", transactionData)

	err := s.d.HandleNotification(notification, transactionData)
	if err != nil {
		fmt.Printf("Error updating transaction: %v\n", err)
		return err
	}

	fmt.Println("Transaction status updated successfully.")
	return nil
}
