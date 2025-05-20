package env

import (
	"github.com/cg917658910/fzkj-wallet/notify-service/lib/log"
	"github.com/joho/godotenv"
)

func init() {
	log.DLogger().Infoln("Loading env file...")

	err := godotenv.Load() // 默认加载当前目录下的 .env 文件
	if err != nil {
		log.DLogger().Fatal("Error loading .env file")
	}
}
