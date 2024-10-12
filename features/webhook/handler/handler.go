package handler

import (
	"fmt"
	"skripsi/features/webhook"

	"github.com/labstack/echo/v4"
)

type WebhookHandler struct {
	s webhook.MidtransNotificationService
}

func New(s webhook.MidtransNotificationService) webhook.MidtransNotificationHandler {
	return &WebhookHandler{
		s: s,
	}
}

func (h *WebhookHandler) HandleNotification() echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("Notification route hit")

		var notification webhook.PaymentNotification
		err := c.Bind(&notification)
		if err != nil {
			return echo.NewHTTPError(400, err.Error())
		}

		err = h.s.HandleNotification(notification)
		if err != nil {
			return echo.NewHTTPError(500, err.Error())
		}

		return c.JSON(200, map[string]string{
			"message": "success"})
	}
}
