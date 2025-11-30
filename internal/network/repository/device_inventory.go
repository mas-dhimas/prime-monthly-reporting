package repository

import (
	"encoding/json"

	"gitlab.com/mas-dhimas/xlsx-prime-monthly-reporting/internal/network/models"
)

func (r *Repository) GetDeviceInventory() (*models.InventoryNode, error) {
	uri := r.configURI.DeviceInventory
	resp, err := r.httpClient.Fetch(uri)
	if err != nil {
		return nil, err
	}

	data := models.InventoryNode{}
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
