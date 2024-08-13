package entity

import (
	"encoding/json"
	"fmt"
)

type BatteryInfo struct {
	BalanceChargeRequest     float64 `json:"balancechargerequest"`
	ChargeCurrentLimit       float64 `json:"chargecurrentlimit"`
	CycleCount               float64 `json:"cyclecount"`
	DischargeCurrentLimit    float64 `json:"dischargecurrentlimit"`
	FullChargeCapacity       float64 `json:"fullchargecapacity"`
	FullChargeCapacityWh     float64 `json:"fullchargecapacitywh"`
	MaximumCellTemperature   float64 `json:"maximumcelltemperature"`
	MaximumCellVoltage       float64 `json:"maximumcellvoltage"`
	MaximumCellVoltageNum    float64 `json:"maximumcellvoltagenum"`
	MaximumModuleCurrent     float64 `json:"maximummodulecurrent"`
	MaximumModuleDcVoltage   float64 `json:"maximummoduledcvoltage"`
	MaximumModuleTemperature float64 `json:"maximummoduletemperature"`
	MinimumCellTemperature   float64 `json:"minimumcelltemperature"`
	MinimumCellVoltage       float64 `json:"minimumcellvoltage"`
	MinimumCellVoltageNum    float64 `json:"minimumcellvoltagenum"`
	MinimumModuleCurrent     float64 `json:"minimummodulecurrent"`
	MinimumModuleDcVoltage   float64 `json:"minimummoduledcvoltage"`
	MinimumModuleTemperature float64 `json:"minimummoduletemperature"`
	NominalModuleDcVoltage   float64 `json:"nominalmoduledcvoltage"`
	RelativeStateOfCharge    float64 `json:"relativestateofcharge"`
	RemainingCapacity        float64 `json:"remainingcapacity"`
	SystemAlarm              float64 `json:"systemalarm"`
	SystemCurrent            float64 `json:"systemcurrent"`
	SystemDcVoltage          float64 `json:"systemdcvoltage"`
	SystemStatus             float64 `json:"systemstatus"`
	SystemTime               float64 `json:"systemtime"`
	SystemWarning            float64 `json:"systemwarning"`
	UsableRemainingCapacity  float64 `json:"usableremainingcapacity"`
}

func ParseBatteryInfo(body []byte) (*BatteryInfo, error) {
	var info BatteryInfo
	err := json.Unmarshal(body, &info)
	if err != nil {
		return nil, fmt.Errorf("unmarshal battery info body: %s", err)
	}
	return &info, nil
}
