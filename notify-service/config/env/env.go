package env

import (
	"github.com/cg917658910/fzkj-wallet/notify-service/lib/log"
	"github.com/joho/godotenv"
)

func init() {
	//var envFile = "./env"
	log.DLogger().Infoln("Loading env file...")
	/* _, filename, _, _ := runtime.Caller(0) // 获取当前文件（config.go）路径
	envPath := path.Dir(filename)          // 获取当前文件目录
	envFile := envPath + "/.env" */
	err := godotenv.Load() // 默认加载当前目录下的 .env 文件
	if err != nil {
		log.DLogger().Fatal("Error loading .env file")
	}
}
