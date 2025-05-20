package db

import (
	"context"
	"fmt"

	"github.com/cg917658910/fzkj-wallet/notify-service/config"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var mysqlDB *gorm.DB

func SetupMysql(ctx context.Context) {
	dsnCnf := &mysqlDriver.Config{
		User:                 config.Configs.MYSQL.User,
		Passwd:               config.Configs.MYSQL.Password,
		Addr:                 config.Configs.MYSQL.Host + fmt.Sprintf(":%s", config.Configs.MYSQL.Port),
		DBName:               config.Configs.MYSQL.DBName,
		AllowNativePasswords: true,
		// WriteTimeout: time.Second * 60,
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSNConfig: dsnCnf,
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	mysqlDB = db
}

func DB() *gorm.DB {
	return mysqlDB
}
