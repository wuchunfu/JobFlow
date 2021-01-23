package loginLogModel

import (
	"gin-vue/middleware/database"
	logger "github.com/sirupsen/logrus"
)

// 用户model
type LoginLog struct {
	Id         int64  `json:"id" gorm:"primary_key; auto_increment; not null"` // 设置字段 id 自增类型。
	Username   string `json:"username" gorm:"type:varchar(50); not null"`      // 用户名
	Ip         string `json:"ip" gorm:"type:varchar(64); not null "`           // ip 地址
	CreateTime string `json:"createTime" gorm:"type:varchar(50); not null"`    // 创建时间
}

func (loginLog *LoginLog) List(page int, pageSize int, username string) ([]LoginLog, int) {
	db := database.GetDB()
	if username != "" {
		db = db.Where("username = ?", username)
	}
	list := make([]LoginLog, 0)
	findErr := db.Model(&loginLog).Offset((page - 1) * pageSize).Limit(pageSize).Order("create_time desc").Find(&list)
	if findErr.Error != nil {
		logger.Error(findErr.Error)
		return nil, -1
	}
	var count int
	countErr := db.Model(&loginLog).Count(&count)
	if countErr.Error != nil {
		logger.Error(countErr.Error)
		return nil, -1
	}
	return list, count
}

// 新增
func (loginLog *LoginLog) Save() {
	db := database.GetDB().Begin()
	err := db.Model(&loginLog).Create(&loginLog)
	if err.Error != nil {
		db.Rollback()
		logger.Error(err.Error)
	}
	db.Commit()
}
