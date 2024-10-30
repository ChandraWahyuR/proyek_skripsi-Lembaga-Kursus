package data

import (
	"fmt"
	"skripsi/features/transaksi"
	transaksiData "skripsi/features/transaksi/data"
	"skripsi/features/webhook"
	"time"

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

func (d *WebhookData) HandleNotification(notification webhook.PaymentNotification, transaction transaksiData.Transaksi) error {
	paymentUpdate := transaksi.UpdateTransaksiStatus{
		ID:     transaction.ID,
		Status: transaction.Status,
	}

	res := d.DB.Begin()
	err := d.DB.Model(&transaksiData.Transaksi{}).Where("id = ?", transaction.ID).Updates(&paymentUpdate).Error
	if err != nil {
		fmt.Printf("Direct update error: %v\n", err)
		return err
	}

	// Tambahan
	historyUpdate := transaksi.UpdateHistoryStatus{
		ID:         transaction.ID,
		Status:     "Active",
		ValidUntil: time.Now().AddDate(0, 3, 0),
	}
	err = d.DB.Model(&transaksiData.TransaksiHistory{}).Where("transaksi_id = ?", transaction.ID).Updates(&historyUpdate).Error
	if err != nil {
		fmt.Printf("Direct update error: %v\n", err)
		return err
	}

	return res.Commit().Error
}
