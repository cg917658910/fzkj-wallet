package config

import (
	"github.com/cg917658910/fzkj-wallet/quote-service/lib/log"
	"github.com/spf13/viper"
)

var Configs *Config

type (
	FiatQuote struct {
		UpdateCacheBaseTime   int64  `mapstructure:"update_cache_base_time" json:"update_cache_base_time" yaml:"update_cache_base_time"`
		UpdateCacheRandomTime int64  `mapstructure:"update_cache_random_time" json:"update_cache_random_time" yaml:"update_cache_random_time"`
		UpdateCacheEnabled    bool   `mapstructure:"update_cache_enabled" json:"update_cache_enabled" yaml:"update_cache_enabled"`
		UpdateCacheHotSymbols string `mapstructure:"update_cache_hot_symbols" json:"update_cache_hot_symbols" yaml:"update_cache_hot_symbols"`
		UpdateCacheHotFiats   string `mapstructure:"update_cache_hot_fiats" json:"update_cache_hot_fiats" yaml:"update_cache_hot_fiats"`
	}
	Config struct {
		FiatQuote FiatQuote `mapstructure:"FiatQuote" json:"FiatQuote" yaml:"FiatQuote"`
	}
)

func init() {
	log.DLogger().Infoln("quotesvc config initializing...")
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
