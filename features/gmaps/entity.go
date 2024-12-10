package gmaps

import "github.com/labstack/echo/v4"

type DirectionsRequest struct {
	Origin      string
	Destination string
}

type DirectionsResponse struct {
	Distance string
	Duration string
	Steps    []string
}

type GmapsHandlerInterface interface {
	GetDirections() echo.HandlerFunc
}

type GmapsDataInterface interface {
	GetDirections(req DirectionsRequest) (DirectionsResponse, error)
}

type GmapsServiceInterface interface {
	GetDirections(req DirectionsRequest) (DirectionsResponse, error)
}
