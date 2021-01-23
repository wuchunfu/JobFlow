package hostModel

import (
	"gin-vue/middleware/database"
	logger "github.com/sirupsen/logrus"
)

// 用户model
type Host struct {
	HostId     int    `json:"hostId" gorm:"primary_key; auto_increment; not null"`
	HostAlias  string `json:"hostAlias" gorm:"type:varchar(100); DEFAULT ''"`
	HostName   string `json:"hostName" gorm:"type:varchar(100); not null"`
	HostPort   int    `json:"hostPort" gorm:"type:int(20); DEFAULT 5921; not null"`
	Remark     string `json:"remark" gorm:"type:varchar(200); not null"`
	CreateTime string `json:"createTime" gorm:"type:varchar(50); not null"`
	UpdateTime string `json:"updateTime" gorm:"type:varchar(50); DEFAULT ''"`
}

func (host *Host) List(page int, pageSize int, hostName string) ([]Host, int) {
	db := database.GetDB()
	if hostName != "" {
		db = db.Where("host_name = ?", hostName)
	}
	list := make([]Host, 0)
	findErr := db.Model(&host).Offset((page - 1) * pageSize).Limit(pageSize).Order("update_time desc").Find(&list)
	if findErr.Error != nil {
		logger.Error(findErr.Error)
		return nil, -1
	}
	var count int
	countErr := db.Model(&host).Count(&count)
	if countErr.Error != nil {
		logger.Error(countErr.Error)
		return nil, -1
	}
	return list, count
}

func (host *Host) AllList() []Host {
	db := database.GetDB()
	list := make([]Host, 0)
	findErr := db.Model(&host).Order("update_time desc").Find(&list)
	if findErr.Error != nil {
		logger.Error(findErr.Error)
		return nil
	}
	return list
}

func (host *Host) Detail(hostId int) *Host {
	db := database.GetDB()
	err := db.Model(&host).Where("host_id = ?", hostId).Find(&host)
	if err.Error != nil {
		logger.Error(err.Error)
		return nil
	}
	return host
}

// 新增
func (host *Host) Save() {
	db := database.GetDB().Begin()
	err := db.Model(&host).Create(&host)
	if err.Error != nil {
		db.Rollback()
		logger.Error(err.Error)
	}
	db.Commit()
}

// 修改
func (host *Host) Update(hostId int, fieldMap map[string]interface{}) {
	db := database.GetDB().Begin()
	err := db.Model(&host).Where("host_id = ?", hostId).Updates(fieldMap)
	if err.Error != nil {
		db.Rollback()
		logger.Error(err.Error)
	}
	db.Commit()
}

// 删除
func (host *Host) Delete(hostId int) {
	db := database.GetDB().Begin()
	err := db.Model(&host).Where("host_id = ?", hostId).Delete(&host)
	if err.Error != nil {
		db.Rollback()
		logger.Error(err.Error)
	}
	db.Commit()
}

// 主机别名是否存在
func (host *Host) IsExistsHostAlias(hostAlias string) int {
	db := database.GetDB()
	var count int
	err := db.Model(&host).Where("host_alias = ?", hostAlias).Count(&count)
	if err.Error != nil {
		logger.Error(err.Error)
		return -1
	}
	return count
}

// 主机名是否存在
func (host *Host) IsExistsHostName(hostName string) int {
	db := database.GetDB()
	var count int
	err := db.Model(&host).Where("host_name = ?", hostName).Count(&count)
	if err.Error != nil {
		logger.Error(err.Error)
		return -1
	}
	return count
}

// 主机是否存在
func (host *Host) IsExistsHost(hostId int) int {
	db := database.GetDB()
	var count int
	err := db.Model(&host).Where("host_id = ?", hostId).Count(&count)
	if err.Error != nil {
		logger.Error(err.Error)
		return -1
	}
	return count
}
