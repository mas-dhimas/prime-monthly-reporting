package repository

import (
	"encoding/json"
	"fmt"

	"gitlab.com/mas-dhimas/xlsx-prime-monthly-reporting/internal/network/models"
)

func (r *Repository) GetDeviceAvailibilityReporting(id string, from, to int) (*models.DeviceAvailabilityReport, error) {
	uri := fmt.Sprintf(r.configURI.DeviceAvailability, id, from, to)
	resp, err := r.httpClient.Fetch(uri)
	if err != nil {
		return nil, err
	}

	data := models.DeviceAvailabilityReport{}
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
