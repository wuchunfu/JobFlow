package database

import (
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/wuchunfu/JobFlow/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"strings"
	"time"
)

var dbPingInterval = 90 * time.Second
var DB *gorm.DB

func InitDB() *gorm.DB {
	setting := config.InitConfig
	dsn := getDbEngineDSN(&setting)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,                    // 慢 SQL 阈值
			LogLevel:      setLogLevel(setting.Log.Level), // Log level
			Colorful:      false,                          // 禁用彩色打印
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   setting.Db.Prefix, // 表名前缀，`User` 的表名应该是 `t_users`
			SingularTable: true,              // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
		PrepareStmt:            true, // 执行任何 SQL 时都创建并缓存预编译语句，可以提高后续的调用速度
		DisableAutomaticPing:   false,
		SkipDefaultTransaction: true, // 对于写操作（创建、更新、删除），为了确保数据的完整性，GORM 会将它们封装在事务内运行。但这会降低性能，你可以在初始化时禁用这种方式
		Logger:                 newLogger,
		AllowGlobalUpdate:      false,
	})
	if err != nil {
		logrus.Errorf("fail to connect database: %v\n", err)
		os.Exit(-1)
	}
	sqlDb, dbErr := db.DB()
	if dbErr != nil {
		logrus.Errorf("fail to connect database: %v\n", dbErr)
		os.Exit(-1)
	}
	// 设置连接池
	// 用于设置连接池中空闲连接的最大数量。
	sqlDb.SetMaxIdleConns(10)
	// 设置打开数据库连接的最大数量
	sqlDb.SetMaxOpenConns(100)

	// * 解决中文字符问题：Error 1366
	db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4")
	go KeepAlivedDb(sqlDb)
	DB = db
	return db
}

func GetDB() *gorm.DB {
	return DB
}

func setLogLevel(logLevel string) logger.LogLevel {
	// 设置日志级别
	level := strings.Replace(strings.ToLower(logLevel), " ", "", -1)
	switch level {
	case "silent":
		return logger.Silent
	case "info":
		return logger.Info
	case "warn":
		return logger.Warn
	case "error":
		return logger.Error
	default:
		return logger.Silent
	}
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

func KeepAlivedDb(engine *sql.DB) {
	t := time.Tick(dbPingInterval)
	var err error
	for {
		<-t
		err = engine.Ping()
		if err != nil {
			logrus.Errorf("database ping error: %v\n", err.Error())
		}
	}
}
