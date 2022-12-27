package main

import "github.com/spf13/viper"

type Config struct {
	WebHook string
	Secret  string
}

var config *Config

func GetConfig() *Config {
	if config == nil {
		config = &Config{
			WebHook: viper.GetString("WEBHOOK"),
			Secret:  viper.GetString("SECRET"),
		}
	}
	return config
}
