package models

type deviceAvailability struct {
	IP         string       `json:"ip_mgmt"`
	NodeName   string       `json:"node_name"`
	Hostname   string       `json:"hostname"`
	Month      string       `json:"month"`
	IcmpPing   availability `json:"icmp_ping"`
	SnmpUptime availability `json:"snmp_uptime"`
}

type availability struct {
	Uptime           int     `json:"uptime"`
	Downtime         int     `json:"downtime"`
	UnknownTime      int     `json:"unknown_time"`
	Ratio            float32 `json:"ratio"`
	DowntimeRatio    float32 `json:"downtime_ratio"`
	UnknownTimeRatio float32 `json:"unknown_time_ratio"`
	Outage           int     `json:"outage"`
}

type DeviceAvailabilityReport struct {
	Data deviceAvailability `json:"data"`
}
