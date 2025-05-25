package config

import (
	"github.com/cg917658910/fzkj-wallet/notify-service/lib/log"
	"github.com/spf13/viper"
)

var Configs *Config

type (
	MYSQL struct {
		Host     string `mapstructure:"host" json:"host" yaml:"host"`
		Port     string `mapstructure:"port" json:"port" yaml:"port"`
		User     string `mapstructure:"user" json:"user" yaml:"user"`
		Password string `mapstructure:"password" json:"password" yaml:"password"`
		DBName   string `mapstructure:"db_name" json:"db_name" yaml:"db_name"`
	}

	Kafka struct {
		Brokers                       string `mapstructure:"brokers" json:"brokers" yaml:"brokers"`
		OrderNofifyTopic              string `mapstructure:"order_notify_topic" json:"order_notify_topic" yaml:"order_notify_topic"`
		OrderNofifyResultTopics       string `mapstructure:"order_notify_result_topics" json:"order_notify_result_topics" yaml:"order_notify_result_topics"`
		OrderNofifyResultDefaultTopic string `mapstructure:"order_notify_result_default_topic" json:"order_notify_result_default_topic" yaml:"order_notify_result_default_topic"`
		OrderNofifyConsumerGroup      string `mapstructure:"order_notify_consumer_group" json:"order_notify_consumer_group" yaml:"order_notify_consumer_group"`
	}
	Redis struct {
		URL                     string `mapstructure:"url" json:"url" yaml:"url"`
		NotifyOrderResultPrefix string `mapstructure:"order_notify_result_prefix" json:"order_notify_result_prefix" yaml:"order_notify_result_prefix"`
	}

	OrderNotify struct {
		OrderNofifyCallerWorkerNum uint `mapstructure:"order_notify_caller_worker_num" json:"order_notify_caller_worker_num" yaml:"order_notify_caller_worker_num"`
		OrderNofifyRetryNum        uint `mapstructure:"order_notify_retry_num" json:"order_notify_retry_num" yaml:"order_notify_retry_num"`
		OrderNofifyRetryDelayTimeS uint `mapstructure:"order_notify_retry_delay_time_s" json:"order_notify_retry_delay_time_s" yaml:"order_notify_retry_delay_time_s"`
	}

	Config struct {
		MYSQL       MYSQL       `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
		Kafka       Kafka       `mapstructure:"kafka" json:"kafka" yaml:"kafka"`
		Redis       Redis       `mapstructure:"redis" json:"redis" yaml:"redis"`
		OrderNotify OrderNotify `mapstructure:"notify" json:"notify" yaml:"notify"`
	}
)

func init() {
	log.DLogger().Infoln("notifysvc config initializing...")
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
