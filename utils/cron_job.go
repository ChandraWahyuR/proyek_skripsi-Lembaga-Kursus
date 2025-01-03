package utils

import (
	"fmt"
	"skripsi/features/transaksi"
	"time"

	"github.com/robfig/cron/v3"
)

func StartCronJob(transaksiService transaksi.TransaksiServiceInterface) {
	c := cron.New()
	// @hourly
	_, err := c.AddFunc("@hourly", func() {
		now := time.Now()
		err := transaksiService.UpdateExpiredTransactions(now)
		if err != nil {
			fmt.Println("Error during cron job:", err)
		} else {
			fmt.Println("Cron job executed successfully")
		}
	})
	if err != nil {
		fmt.Println("Failed to initialize cron job:", err)
		return
	}
	c.Start()
	fmt.Println("Cron job started")
}
