package main

import (
	"flag"
	"gok-pi/battery/api-client"
	"gok-pi/battery/discharger"
	"gok-pi/internal/config"
	"gok-pi/internal/lib/logger"
	"gok-pi/internal/lib/sl"
	"log/slog"
	"sync"
)

func main() {

	configPath := flag.String("conf", "config.yml", "path to config file")
	logPath := flag.String("log", "/var/log", "path to log file directory")
	flag.Parse()

	conf := config.MustLoad(*configPath)
	lg := logger.SetupLogger(conf.Env, *logPath)

	lg.Info("starting gok-pi", slog.String("config", *configPath), slog.String("env", conf.Env))
	lg.Debug("debug messages enabled")
	// filter enabled batteries
	var batteries []config.BatteryConfig
	for _, b := range conf.Batteries {
		if b.Enabled {
			batteries = append(batteries, b)
		}
	}
	lg.With(
		slog.Int("batteries", len(batteries)),
	).Info("loaded batteries config")

	if len(batteries) == 0 {
		lg.Warn("no batteries enabled")
		return
	}

	var wg sync.WaitGroup

	for _, b := range batteries {
		wg.Add(1)
		go func(workerId string) {
			defer wg.Done()

			log := lg.With(slog.String("battery", workerId))
			api := apiclient.New(b.Url, b.Token, log)

			worker, err := discharger.New(conf.StartTime, conf.StopTime, b.CapacityLimit, b.PowerLimit, api, log)
			if err != nil {
				log.Error("creating discharge worker", sl.Err(err))
			}
			err = worker.Run()
			if err != nil {
				log.Error("running discharge worker", sl.Err(err))
			}
			log.Info("discharge worker stopped")
		}(b.Name)
	}
	wg.Wait()

	lg.Info("gok-pi stopped")
}
