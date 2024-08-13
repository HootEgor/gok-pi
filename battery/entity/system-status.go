package entity

import (
	"encoding/json"
	"fmt"
)

type SystemStatus struct {
	ApparentOutput            float64     `json:"Apparent_output"`
	BackupBuffer              string      `json:"BackupBuffer"`
	BatteryCharging           bool        `json:"BatteryCharging"`
	BatteryDischarging        bool        `json:"BatteryDischarging"`
	ConsumptionAvg            float64     `json:"Consumption_Avg"`
	ConsumptionW              float64     `json:"Consumption_W"`
	Fac                       float64     `json:"Fac"`
	FlowConsumptionBattery    bool        `json:"FlowConsumptionBattery"`
	FlowConsumptionGrid       bool        `json:"FlowConsumptionGrid"`
	FlowConsumptionProduction bool        `json:"FlowConsumptionProduction"`
	FlowGridBattery           bool        `json:"FlowGridBattery"`
	FlowProductionBattery     bool        `json:"FlowProductionBattery"`
	FlowProductionGrid        bool        `json:"FlowProductionGrid"`
	GridFeedInW               float64     `json:"GridFeedIn_W"`
	IsSystemInstalled         float64     `json:"IsSystemInstalled"`
	OperatingMode             string      `json:"OperatingMode"`
	PacTotalW                 float64     `json:"Pac_total_W"`
	ProductionW               float64     `json:"Production_W"`
	RSOC                      float64     `json:"RSOC"`
	RemainingCapacityWh       float64     `json:"RemainingCapacity_Wh"`
	Sac1                      float64     `json:"Sac1"`
	Sac2                      interface{} `json:"Sac2"`
	Sac3                      interface{} `json:"Sac3"`
	SystemStatus              string      `json:"SystemStatus"`
	Timestamp                 string      `json:"Timestamp"`
	USOC                      float64     `json:"USOC"`
	Uac                       float64     `json:"Uac"`
	Ubat                      float64     `json:"Ubat"`
	DischargeNotAllowed       bool        `json:"dischargeNotAllowed"`
	GeneratorAutostart        bool        `json:"generator_autostart"`
}

func ParseSystemStatus(body []byte) (*SystemStatus, error) {
	var status SystemStatus
	err := json.Unmarshal(body, &status)
	if err != nil {
		return nil, fmt.Errorf("unmarshal status body: %s", err)
	}
	return &status, nil
}
