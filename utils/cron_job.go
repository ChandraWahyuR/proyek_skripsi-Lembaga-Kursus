package utils

import (
	"fmt"
	"log"
	"net/http"
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

func CronJobMinute() {
	c := cron.New()
	_, err := c.AddFunc("*/12 * * * *", func() {
		resp, err := http.Get("https://skripsi-245802795341.asia-southeast2.run.app/ping")
		if err != nil {
			log.Printf("Ping failed: %v", err)
			return
		}
		defer resp.Body.Close()
		log.Println("Ping successful:", resp.Status)
	})
	if err != nil {
		log.Printf("Error adding cron job: %v", err)
		return
	}

	c.Start()
}
