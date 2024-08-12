package main

import (
	"flag"
	"gok-pi/battery/api-client"
	"gok-pi/battery/discharger"
	"gok-pi/internal/config"
	"gok-pi/internal/lib/logger"
	"gok-pi/internal/lib/sl"
	"log/slog"
)

func main() {

	configPath := flag.String("conf", "config.yml", "path to config file")
	logPath := flag.String("log", "/var/log", "path to log file directory")
	flag.Parse()

	conf := config.MustLoad(*configPath)
	lg := logger.SetupLogger(conf.Env, *logPath)

	lg.Info("starting gok-pi", slog.String("config", *configPath), slog.String("env", conf.Env))
	lg.Debug("debug messages enabled")

	api := apiclient.New(conf.Endpoint.Url, conf.Endpoint.Token, lg)

	worker, err := discharger.New(conf.StartTime, conf.StopTime, conf.BatteryLimit, api, lg)
	if err != nil {
		lg.Error("creating discharge worker", sl.Err(err))
	}
	err = worker.Run()
	if err != nil {
		lg.Error("running discharge worker", sl.Err(err))
	}
}
