package helper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type GeocodeResponse struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

func GetCoordinates(address string) (float64, float64, error) {
	baseURL := "https://nominatim.openstreetmap.org/search"
	params := url.Values{}
	params.Add("q", address)
	params.Add("format", "json")

	resp, err := http.Get(fmt.Sprintf("%s?%s", baseURL, params.Encode()))
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	var results []GeocodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return 0, 0, err
	}

	if len(results) == 0 {
		return 0, 0, fmt.Errorf("no coordinates found for address: %s", address)
	}

	lat, _ := strconv.ParseFloat(results[0].Lat, 64)
	lon, _ := strconv.ParseFloat(results[0].Lon, 64)

	return lat, lon, nil
}
