package mysql

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/icowan/shalom/src/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

func NewDb(logger log.Logger, cf *config.Config) (*gorm.DB, error) {
	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local&timeout=20m&collation=utf8mb4_unicode_ci",
		cf.GetString("mysql", "user"), cf.GetString("mysql", "password"),
		cf.GetString("mysql", "host"), cf.GetString("mysql", "port"),
		cf.GetString("mysql", "database"))

	db, err := gorm.Open("mysql", dbUrl)
	if err != nil {
		return nil, err
	}

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.DB().SetConnMaxLifetime(time.Hour)
	db.LogMode(cf.GetBool("server", "debug"))
	//db.Raw("SET sql_mode = '';")
	//var c interface{}
	//if err = db.Raw("select @@version").Scan(c).Error; err != nil {
	//	_ = logger.Log("version", err.Error())
	//}
	//

	if err = db.Raw("set global sql_mode='STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION';").Error; err != nil {
		_ = logger.Log("db.Raw", err.Error())
	}
	if err = db.Raw("set session sql_mode='STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION';").Error; err != nil {
		_ = logger.Log("db.Raw", err.Error())
	}
	if err = db.DB().Ping(); err != nil {
		_ = logger.Log("db", "ping", "err", err)
	}
	
	_ = logger.Log("mysql", "connect", "success", true)

	return db, nil
}
