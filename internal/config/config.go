package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

type Config struct {
	Env      string `yaml:"env" env-default:"local" env-required:"true"`
	Endpoint struct {
		Url   string `yaml:"url" env-default:"https://example.battery/api"`
		Token string `yaml:"token" env-default:"auth-token"`
	} `yaml:"endpoint"`
	StartTime    string `yaml:"start_time" env-default:"18:00"`
	StopTime     string `yaml:"stop_time" env-default:"22:00"`
	BatteryLimit int    `yaml:"battery_limit" env-default:"20"`
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
