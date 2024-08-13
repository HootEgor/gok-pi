package entity

import (
	"encoding/json"
	"fmt"
)

type SystemStatus struct {
	ApparentOutput            int         `json:"Apparent_output"`
	BackupBuffer              string      `json:"BackupBuffer"`
	BatteryCharging           bool        `json:"BatteryCharging"`
	BatteryDischarging        bool        `json:"BatteryDischarging"`
	ConsumptionAvg            int         `json:"Consumption_Avg"`
	ConsumptionW              int         `json:"Consumption_W"`
	Fac                       float64     `json:"Fac"`
	FlowConsumptionBattery    bool        `json:"FlowConsumptionBattery"`
	FlowConsumptionGrid       bool        `json:"FlowConsumptionGrid"`
	FlowConsumptionProduction bool        `json:"FlowConsumptionProduction"`
	FlowGridBattery           bool        `json:"FlowGridBattery"`
	FlowProductionBattery     bool        `json:"FlowProductionBattery"`
	FlowProductionGrid        bool        `json:"FlowProductionGrid"`
	GridFeedInW               int         `json:"GridFeedIn_W"`
	IsSystemInstalled         int         `json:"IsSystemInstalled"`
	OperatingMode             string      `json:"OperatingMode"`
	PacTotalW                 int         `json:"Pac_total_W"`
	ProductionW               int         `json:"Production_W"`
	RSOC                      int         `json:"RSOC"`
	RemainingCapacityWh       int         `json:"RemainingCapacity_Wh"`
	Sac1                      int         `json:"Sac1"`
	Sac2                      interface{} `json:"Sac2"`
	Sac3                      interface{} `json:"Sac3"`
	SystemStatus              string      `json:"SystemStatus"`
	Timestamp                 string      `json:"Timestamp"`
	USOC                      int         `json:"USOC"`
	Uac                       int         `json:"Uac"`
	Ubat                      int         `json:"Ubat"`
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
