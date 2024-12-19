package handler

import (
	"net/http"
	"skripsi/features/gmaps"
	"skripsi/helper"

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
		const destination = "-7.484362621456366,108.78844442023427"

		if origin == "" {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Masukkan alamat asal", nil))
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
