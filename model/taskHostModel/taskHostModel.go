package taskHostModel

import (
	logger "github.com/sirupsen/logrus"
	"github.com/wuchunfu/JobFlow/middleware/database"
)

type TaskHost struct {
	Id        int    `json:"id" gorm:"type:int(20); primary_key; auto_increment; not null"`
	TaskId    int    `json:"taskId" gorm:"type:int(20); not null; index:IDX_task_host_task_id"`
	HostId    int    `json:"hostId" gorm:"type:int(20); not null; index:IDX_task_host_host_id"`
	HostName  string `json:"hostName" gorm:"-"`
	HostPort  int    `json:"hostPort" gorm:"-"`
	HostAlias string `json:"hostAlias" gorm:"-"`
}

// 新增
func (taskHost *TaskHost) Save(taskId int, hostIds []int) {
	db := database.GetDB().Begin()
	//taskHost.Delete(taskId)

	//taskHosts := make([]TaskHost, len(hostIds))
	//for i, value := range hostIds {
	//	taskHosts[i].TaskId = taskId
	//	taskHosts[i].HostId = value
	//}
	//taskHosts := make([]TaskHost, len(hostIds))
	for _, value := range hostIds {
		taskHost.TaskId = taskId
		taskHost.HostId = value
	}

	err := db.Model(&taskHost).Create(&taskHost)
	if err.Error != nil {
		db.Rollback()
		logger.Error(err.Error)
	}
	db.Commit()
}

// 删除
func (taskHost *TaskHost) Delete(id int) {
	db := database.GetDB().Begin()
	err := db.Model(&taskHost).Where("id = ?", id).Delete(&taskHost)
	if err.Error != nil {
		db.Rollback()
		logger.Error(err.Error)
	}
	db.Commit()
}

func (taskHost *TaskHost) GetTaskIdsByHostId(hostId int) ([]interface{}, error) {
	db := database.GetDB()
	list := make([]TaskHost, 0)
	err := db.Model(&taskHost).Where("host_id = ?", hostId).Select("task_id").Find(&list)
	if err.Error != nil {
		return nil, err.Error
	}
	taskIds := make([]interface{}, len(list))
	for i, value := range list {
		taskIds[i] = value.TaskId
	}
	return taskIds, err.Error
}

func (taskHost *TaskHost) GetHostIdsByTaskId(taskId int) []TaskHost {
	db := database.GetDB()
	list := make([]TaskHost, 0)
	err := db.Raw("select th.id,th.host_id,h.host_alias,h.host_name,h.host_port from gin_task_host as th left join gin_host as h on th.host_id = h.host_id where th.task_id = ?", taskId).Find(&list)
	if err.Error != nil {
		logger.Error(err.Error)
		return nil
	}
	return list
}
