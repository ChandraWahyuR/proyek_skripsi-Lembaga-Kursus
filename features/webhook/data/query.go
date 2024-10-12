package data

import (
	"fmt"
	"skripsi/features/transaksi"
	"skripsi/features/webhook"

	"gorm.io/gorm"
)

type WebhookData struct {
	DB *gorm.DB
}

func New(db *gorm.DB) webhook.MidtransNotificationData {
	return &WebhookData{
		DB: db,
	}
}

func (d *WebhookData) HandleNotification(notification webhook.PaymentNotification, transaction transaksi.Transaksi) error {
	paymentUpdate := transaksi.UpdateTransaksiStatus{
		ID:     transaction.ID,
		Status: transaction.Status,
	}

	res := d.DB.Begin()
	err := d.DB.Model(&transaksi.Transaksi{}).Where("id = ?", transaction.ID).Updates(&paymentUpdate).Error
	if err != nil {
		fmt.Printf("Direct update error: %v\n", err)
		return err
	}

	return res.Commit().Error
}
