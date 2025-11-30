package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type inventoryNodeResponse struct {
	ID                       primitive.ObjectID  `json:"id"`
	NodeName                 string              `json:"node_name"`
	IPMgmt                   string              `json:"ip_mgmt"`
	ProbeID                  primitive.ObjectID  `json:"probe_id,omitempty"`
	ProbeInstallationType    string              `json:"probe_installation_type,omitempty"`
	ProbeKey                 string              `json:"probe_key,omitempty"`
	AgentName                string              `json:"probe_name,omitempty"`
	Notes                    string              `json:"notes,omitempty"`
	Threshold                float64             `json:"threshold"`
	TrafficCounterPreference string              `json:"traffic_counter_preference"`
	Modules                  map[string][]string `json:"modules"`
	Username                 string              `json:"username,omitempty"`
	Password                 string              `json:"password,omitempty"`
	CreatedAt                time.Time           `json:"created_at,omitempty"`
	UpdatedAt                time.Time           `json:"updated_at,omitempty"`
}

type paginatedNodeResponse struct {
	Nodes     []inventoryNodeResponse `json:"nodes"`
	TotalData int64                   `json:"total_data"`
}

type InventoryNode struct {
	Data paginatedNodeResponse `json:"data"`
}
