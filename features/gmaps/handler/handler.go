package handler

import (
	"net/http"
	"skripsi/features/gmaps"

	"github.com/labstack/echo/v4"
)

type GmapsHandler struct {
	s gmaps.GmapsServiceInterface
}

func New(u gmaps.GmapsServiceInterface) gmaps.GmapsHandlerInterface {
	return &GmapsHandler{
		s: u,
	}
}

func (h *GmapsHandler) GetDirections() echo.HandlerFunc {
	return func(c echo.Context) error {
		origin := c.QueryParam("origin")
		destination := c.QueryParam("destination")

		if origin == "" || destination == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "origin and destination are required"})
		}

		request := gmaps.DirectionsRequest{
			Origin:      origin,
			Destination: destination,
		}

		response, err := h.s.GetDirections(request)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}

		return c.JSON(http.StatusOK, response)
	}
}
