package database

import (
	"fmt"
	"gin-vue/config"
	"github.com/jinzhu/gorm"
	logger "github.com/sirupsen/logrus"
	"strings"
	"time"
)

var dbPingInterval = 90 * time.Second
var DB *gorm.DB

func InitDB() *gorm.DB {
	setting := config.InitConfig
	dsn := getDbEngineDSN(&setting)
	db, err := gorm.Open(setting.Db.Engine, dsn)
	if err != nil {
		panic("fail to connect database, err:" + err.Error())
	}
	// 即用复数形式
	db.SingularTable(true)
	// 设置连接池
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.LogMode(true)
	// 为表名添加前缀
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return setting.Db.Prefix + defaultTableName
	}
	// * 解决中文字符问题：Error 1366
	db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4")
	go keepDbAlived(db)
	DB = db
	return db
}

func GetDB() *gorm.DB {
	return DB
}

// 获取数据库引擎DSN  mysql,postgres
func getDbEngineDSN(setting *config.Server) string {
	engine := strings.ToLower(setting.Db.Engine)
	dsn := ""
	switch engine {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&allowNativePasswords=true&parseTime=True&loc=Local",
			setting.Db.Username,
			setting.Db.Password,
			setting.Db.Host,
			setting.Db.Port,
			setting.Db.Database,
			setting.Db.Charset)
	case "postgres":
		dsn = fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
			setting.Db.Host,
			setting.Db.Port,
			setting.Db.Database,
			setting.Db.Username,
			setting.Db.Password)
	}
	return dsn
}

func keepDbAlived(engine *gorm.DB) {
	t := time.Tick(dbPingInterval)
	var err error
	for {
		<-t
		err = engine.DB().Ping()
		if err != nil {
			logger.Infof("database ping: %s", err)
		}
	}
}
