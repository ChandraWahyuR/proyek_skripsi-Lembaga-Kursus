package data

import (
	"encoding/json"
	"fmt"
	"net/http"
	"skripsi/features/gmaps"
)

type GmapsData struct {
	apiKey string
}

func New(apiKey string) gmaps.GmapsDataInterface {
	return &GmapsData{
		apiKey: apiKey,
	}
}

func (r *GmapsData) GetDirections(req gmaps.DirectionsRequest) (gmaps.DirectionsResponse, error) {
	// URL untuk API Directions Google Maps
	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/directions/json?origin=%s&destination=%s&key=%s",
		req.Origin, req.Destination, r.apiKey)

	// Melakukan HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return gmaps.DirectionsResponse{}, err
	}
	defer resp.Body.Close()

	// Validasi response status
	if resp.StatusCode != http.StatusOK {
		return gmaps.DirectionsResponse{}, fmt.Errorf("error from Google Maps API: %s", resp.Status)
	}

	// Parsing JSON response dari Google Maps
	var result struct {
		Routes []struct {
			Legs []struct {
				Distance struct {
					Text string `json:"text"`
				} `json:"distance"`
				Duration struct {
					Text string `json:"text"`
				} `json:"duration"`
				Steps []struct {
					HTMLInstructions string `json:"html_instructions"`
				} `json:"steps"`
			} `json:"legs"`
		} `json:"routes"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return gmaps.DirectionsResponse{}, err
	}

	// Validasi apakah ada route di response
	if len(result.Routes) == 0 || len(result.Routes[0].Legs) == 0 {
		return gmaps.DirectionsResponse{}, fmt.Errorf("no routes found")
	}

	leg := result.Routes[0].Legs[0]
	steps := make([]string, len(leg.Steps))
	for i, step := range leg.Steps {
		steps[i] = step.HTMLInstructions
	}

	return gmaps.DirectionsResponse{
		Distance: leg.Distance.Text,
		Duration: leg.Duration.Text,
		Steps:    steps,
	}, nil
}
