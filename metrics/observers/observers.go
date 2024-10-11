package observers

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"strconv"
)

var socGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "battery",
	Name:      "RSoC",
	Help:      "Relative state of charge in percent",
}, []string{"name"})

func UpdateSoC(name string, value float64) {
	socGauge.WithLabelValues(name).Set(value)
}

var uSocGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "battery",
	Name:      "USoC",
	Help:      "User state of charge in percent",
}, []string{"name"})

func UpdateUSoC(name string, value float64) {
	uSocGauge.WithLabelValues(name).Set(value)
}

var capacityGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "battery",
	Name:      "RemainingCapacity_W",
	Help:      "Remaining capacity based on RSoC",
}, []string{"name"})

func UpdateCapacity(name string, value float64) {
	capacityGauge.WithLabelValues(name).Set(value)
}

var consumptionGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "battery",
	Name:      "Consumption_W",
	Help:      "House consumption in Watts, direct measurement",
}, []string{"name"})

func UpdateConsumption(name string, value float64) {
	consumptionGauge.WithLabelValues(name).Set(value)
}

var pacGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "battery",
	Name:      "Pac_total_W",
	Help:      "AC Power: greater than zero - discharging, less than zero - charging in Watts",
}, []string{"name"})

func UpdatePac(name string, value float64) {
	pacGauge.WithLabelValues(name).Set(value)
}

var dischargeStateGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "battery",
	Name:      "BatteryDischarging",
	Help:      "Discharge status: 1 - discharging, 0 - not discharging",
}, []string{"name"})

func UpdateDischargeState(name string, state bool) {
	if state {
		dischargeStateGauge.WithLabelValues(name).Set(1.0)
	} else {
		dischargeStateGauge.WithLabelValues(name).Set(0.0)
	}
}

var opModeGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "battery",
	Name:      "BatteryOperatingMode",
	Help:      "Operating mode: 1 - manual, 2 - auto",
}, []string{"name"})

func UpdateOpMode(name string, value string) {
	state, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return
	}
	opModeGauge.WithLabelValues(name).Set(state)
}
