package main

import "github.com/spf13/viper"

type Config struct {
	Logger LoggerConf
	HTTP   HTTPConf
}

type LoggerConf struct {
	Level string
}

func NewConfig() Config {
	return Config{}
}

type HTTPConf struct {
	Port int
}

func ReadConfig(path string) (*Config, error) {
	v := viper.New()
	v.AutomaticEnv()
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	config := NewConfig()
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
