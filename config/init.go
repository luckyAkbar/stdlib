// Package config is used to read the config file
package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")
	viper.AddConfigPath("../../")
	viper.AddConfigPath("../../../")

	err := viper.ReadInConfig()
	if err != nil {
		logrus.Error(err)
		panic("failed to read config file")
	}
}
