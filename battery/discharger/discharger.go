package discharger

import (
	"fmt"
	"gok-pi/battery/entity"
	"gok-pi/internal/lib/sl"
	"gok-pi/internal/lib/timer"
	"log/slog"
	"time"
)

type Client interface {
	Status() (*entity.SystemStatus, error)
	StartDischarge(power int) error
	StopDischarge() error
}

type Discharge struct {
	startTime    string
	stopTime     string
	batteryLimit float64
	powerLimit   int
	client       Client
	log          *slog.Logger
}

func New(startTime, stopTime string, batteryLimit, powerLimit int, client Client, log *slog.Logger) (*Discharge, error) {
	return &Discharge{
		startTime:    startTime,
		stopTime:     stopTime,
		batteryLimit: float64(batteryLimit),
		powerLimit:   powerLimit,
		client:       client,
		log:          log.With(sl.Module("battery.discharge")),
	}, nil
}

func (d *Discharge) Run() error {
	for {
		// Calculate the start and stop times for today
		startTime, err := timer.ParseTime(d.startTime)
		if err != nil {
			return fmt.Errorf("parsing start time: %w", err)
		}
		stopTime, err := timer.ParseTime(d.stopTime)
		if err != nil {
			return fmt.Errorf("parsing stop time: %w", err)
		}
		if startTime.After(stopTime) {
			stopTime = stopTime.Add(24 * time.Hour)
		}
		now := time.Now()
		d.log.With(
			slog.String("start_time", startTime.Format(time.DateTime)),
			slog.String("stop_time", stopTime.Format(time.DateTime)),
			slog.String("now", now.Format(time.DateTime)),
			slog.Float64("limit", d.batteryLimit),
		).Info("next cycle")

		// If start time has passed for today, schedule for the next day
		if now.After(stopTime) {
			startTime = startTime.Add(24 * time.Hour)
		}

		startTimer := time.NewTimer(startTime.Sub(now))

		d.log.With(slog.Time("start_time", startTime)).Info("waiting until start time")
		<-startTimer.C

		// Check the battery status
		d.log.With(slog.Float64("limit", d.batteryLimit)).Info("starting battery discharge process...")
		status, err := d.client.Status()
		if err != nil {
			d.log.With(sl.Err(err)).Error("checking battery status")
			time.Sleep(1 * time.Minute)
			continue
		}
		d.log.With(
			slog.Float64("remaining capacity", status.RemainingCapacityWh),
			slog.Float64("consumption", status.ConsumptionW),
			slog.Bool("discharge", status.BatteryDischarging),
		).Info("status")

		if status.RemainingCapacityWh > d.batteryLimit {
			//err = d.client.StartDischarge()
			//if err != nil {
			//	d.log.With(sl.Err(err)).Error("starting discharge")
			//}

			// Start monitoring battery status during discharge
			d.monitorState(stopTime)

		} else {
			d.log.Info("battery level is below the limit, no discharge needed")
		}

		d.log.Info("waiting for the next cycle...")
		time.Sleep(24*time.Hour - time.Now().Sub(startTime))
	}
}

func (d *Discharge) monitorState(stopTime time.Time) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	stopTimer := time.NewTimer(stopTime.Sub(time.Now()))

	for {
		select {
		case <-ticker.C:
			status, err := d.client.Status()
			if err != nil {
				d.log.With(sl.Err(err)).Error("checking battery status")
				continue
			}
			log := d.log.With(
				slog.Float64("remaining capacity", status.RemainingCapacityWh),
				slog.Float64("consumption", status.ConsumptionW),
				slog.Bool("discharge", status.BatteryDischarging),
			)

			if status.RemainingCapacityWh <= d.batteryLimit && status.BatteryDischarging {
				log.Info("battery level reached the limit, stopping discharge")
				err = d.client.StopDischarge()
				if err != nil {
					d.log.With(sl.Err(err)).Error("stopping discharge")
					continue
				}
				return
			}

			if status.BatteryDischarging == false {

				dischargePower := d.calculateRate(status.RemainingCapacityWh, stopTime)
				log = log.With(slog.Int("discharge_power", dischargePower))
				if dischargePower == 0 {
					log.Warn("discharge power is zero, stopping discharge")
					err = d.client.StopDischarge()
					if err != nil {
						d.log.With(sl.Err(err)).Error("stopping discharge")
						continue
					}
					return
				}

				log.Info("starting discharge")
				err = d.client.StartDischarge(dischargePower)
				if err != nil {
					d.log.With(sl.Err(err)).Error("starting discharge")
				}

			} else {
				log.Info("discharge is running")
			}

		case <-stopTimer.C:
			d.log.Info("stop time reached, stopping discharge")
			err := d.client.StopDischarge()
			if err != nil {
				d.log.With(sl.Err(err)).Error("stopping discharge")
			}
			return
		}
	}
}

// calculate discharge rate as Wh/h
func (d *Discharge) calculateRate(capacity float64, stopTime time.Time) int {
	estimate := capacity - d.batteryLimit
	if estimate <= 0 {
		return 0
	}
	remainingTime := stopTime.Sub(time.Now())
	if remainingTime <= 0 {
		return 0
	}
	rate := estimate / remainingTime.Hours()
	if rate > float64(d.powerLimit) {
		return d.powerLimit
	}
	return int(rate)
}
