package config

import (
	"github.com/cg917658910/fzkj-wallet/notify-service/lib/log"
	"github.com/spf13/viper"
)

var Configs *Config

type MYSQL struct {
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Port     string `mapstructure:"port" json:"port" yaml:"port"`
	User     string `mapstructure:"user" json:"user" yaml:"user"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	DBName   string `mapstructure:"db_name" json:"db_name" yaml:"db_name"`
}

type Kafka struct {
	Brokers                  string `mapstructure:"brokers" json:"brokers" yaml:"brokers"`
	OrderNofifyTopic         string `mapstructure:"order_notify_topic" json:"order_notify_topic" yaml:"order_notify_topic"`
	OrderNofifyResultTopic   string `mapstructure:"order_notify_result_topic" json:"order_notify_result_topic" yaml:"order_notify_result_topic"`
	OrderNofifyConsumerGroup string `mapstructure:"order_notify_consumer_group" json:"order_notify_consumer_group" yaml:"order_notify_consumer_group"`
}

type Config struct {
	MYSQL MYSQL `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Kafka Kafka `mapstructure:"kafka" json:"kafka" yaml:"kafka"`
}

func init() {
	log.DLogger().Infoln("notifysvc config initializing...")

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		panic("Config Read failed: " + err.Error())
	}
	err := viper.Unmarshal(&Configs)
	if err != nil {
		panic("Config decode failed: " + err.Error())
	}
}
