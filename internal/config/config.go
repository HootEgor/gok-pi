package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

type Config struct {
	Env       string          `yaml:"env" env-default:"local" env-required:"true"`
	StartTime string          `yaml:"start_time" env-default:"18:00"`
	StopTime  string          `yaml:"stop_time" env-default:"22:00"`
	Metrics   MetricsServer   `yaml:"metrics"`
	Batteries []BatteryConfig `yaml:"batteries"`
}

type BatteryConfig struct {
	Name          string `yaml:"name" env-default:"battery1"`
	Url           string `yaml:"url" env-default:"https://example.battery/api"`
	Token         string `yaml:"token" env-default:"auth-token"`
	Enabled       bool   `yaml:"enabled" env-default:"true"`
	CapacityLimit int    `yaml:"capacity_limit" env-default:"20000"`
	PowerLimit    int    `yaml:"power_limit" env-default:"1000"`
	SocLimit      int    `yaml:"soc_limit" env-default:"50"`
}

type MetricsServer struct {
	Enabled bool   `yaml:"enabled" env-default:"false"`
	Bind    string `yaml:"bind" env-default:"0.0.0.0"`
	Port    string `yaml:"port" env-default:"5001"`
}

var instance *Config
var once sync.Once

func MustLoad(path string) *Config {
	var err error
	once.Do(func() {
		instance = &Config{}
		if err = cleanenv.ReadConfig(path, instance); err != nil {
			desc, _ := cleanenv.GetDescription(instance, nil)
			err = fmt.Errorf("%s; %s", err, desc)
			instance = nil
			log.Fatal(err)
		}
	})
	return instance
}
