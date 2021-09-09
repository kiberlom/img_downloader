package config

import "github.com/spf13/viper"

type Config struct {
	C *viper.Viper
}

func NewConfig() *Config {
	s := viper.New()
	s.AutomaticEnv()

	return &Config{C: s}

}
