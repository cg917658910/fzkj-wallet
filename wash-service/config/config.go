package config

import (
	"github.com/cg917658910/fzkj-wallet/wash-service/lib/log"
	"github.com/spf13/viper"
)

var Configs *Config

type (
	Config struct {
	}
)

func init() {
	log.DLogger().Infoln("washsvc config initializing...")
	var confPath = "."
	//_, filename, _, _ := runtime.Caller(0) // 获取当前文件（config.go）路径
	//confPath := path.Dir(filename)         // 获取当前文件目录
	viper.SetConfigName("config")
	viper.AddConfigPath(confPath)
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")
	viper.WatchConfig()

	if err := viper.ReadInConfig(); err != nil {
		panic("Config Read failed: " + err.Error())
	}
	err := viper.Unmarshal(&Configs)
	if err != nil {
		panic("Config decode failed: " + err.Error())
	}
}
