package taskLogModel

import (
	"github.com/sirupsen/logrus"
	"github.com/wuchunfu/JobFlow/common"
	"github.com/wuchunfu/JobFlow/middleware/database"
	"github.com/wuchunfu/JobFlow/utils/datetimeUtils"
	"time"
)

type TaskType int8

// 任务执行日志
type TaskLog struct {
	Id             int64  `json:"id" gorm:"type:bigint(20); primary_key; auto_increment; not null"`
	TaskId         int    `json:"taskId" gorm:"type:int(20); not null; index:IDX_task_log_task_id; default 0"` // 任务id
	TaskName       string `json:"taskName" gorm:"type:varchar(100); not null"`                                 // 任务名称
	CronExpression string `json:"cronExpression" gorm:"type:varchar(100); not null"`                           // crontab
	Protocol       int8   `json:"protocol" gorm:"type:int(11); not null; index:IDX_task_log_protocol"`         // 协议 1:http 2:RPC
	Command        string `json:"command" gorm:"type:varchar(300); not null"`                                  // URL地址或shell命令
	Timeout        int    `json:"timeout" gorm:"type:int(11); not null; default 0"`                            // 任务执行超时时间(单位秒),0不限制
	RetryTimes     int8   `json:"retryTimes" gorm:"type:int(11); not null; default 0"`                         // 任务重试次数
	HostName       string `json:"hostName" gorm:"type:varchar(200); not null; default ''"`                     // RPC主机名，逗号分隔
	StartTime      string `json:"startTime" gorm:"type:varchar(50); default null"`                             // 开始执行时间
	EndTime        string `json:"endTime" gorm:"type:varchar(50); default null"`                               // 执行完成（失败）时间
	Status         int8   `json:"status" gorm:"type:int(11); not null; index:IDX_task_log_status; default 1"`  // 状态 0:执行失败 1:执行中  2:执行完毕 3:任务取消(上次任务未执行完成) 4:异步执行
	Result         string `json:"result" gorm:"type:text; not null"`                                           // 执行结果
	TotalTime      int    `json:"totalTime" gorm:"-"`                                                          // 执行总时长
}

func (taskLog *TaskLog) Create() int64 {
	db := database.GetDB()
	err := db.Model(&taskLog).Create(&taskLog)
	if err.Error != nil {
		logrus.Error(err.Error)
	}
	return taskLog.Id
}

// 更新
func (taskLog *TaskLog) Update(id int64, fieldMap map[string]interface{}) int64 {
	db := database.GetDB()
	err := db.Model(&taskLog).Where("id = ?", id).Updates(fieldMap)
	if err.Error != nil {
		logrus.Error(err.Error)
	}
	return taskLog.Id
}

func (taskLog *TaskLog) List(page int, pageSize int, taskName string) ([]TaskLog, int64) {
	db := database.GetDB()
	if taskName != "" {
		db = db.Where("task_name = ?", taskName)
	}
	list := make([]TaskLog, 0)
	err := db.Model(&taskLog).Offset((page - 1) * pageSize).Limit(pageSize).Order("id desc").Find(&list)
	if len(list) > 0 {
		for i, item := range list {
			endTime := item.EndTime
			if item.Status == common.Running {
				//endTime = time.Now()
				endTime = datetimeUtils.FormatDateTime()
			}
			dateTime2, _ := time.Parse("2006-01-02 15:04:05", endTime)
			dateTime1, _ := time.Parse("2006-01-02 15:04:05", item.StartTime)
			//execSeconds := endTime.Sub(item.StartTime).Seconds()
			execSeconds := dateTime2.Sub(dateTime1).Seconds()
			list[i].TotalTime = int(execSeconds)
		}
	}
	if err.Error != nil {
		logrus.Error(err.Error)
		return nil, -1
	}
	var count int64
	countErr := db.Model(&taskLog).Count(&count)
	if countErr.Error != nil {
		logrus.Error(countErr.Error)
		return nil, -1
	}
	return list, count
}

// 清空表
func (taskLog *TaskLog) Clear() {
	db := database.GetDB()
	db.Model(&taskLog).Where("1=1").Delete(&taskLog)
}

// 删除N个月前的日志
func (taskLog *TaskLog) Remove(id int) {
	now := time.Now().AddDate(0, -id, 0)
	db := database.GetDB()
	db.Model(&taskLog).Where("start_time <= ?", now.Format("2006-01-02 15:04:05")).Delete(&taskLog)
}

func (taskLog *TaskLog) Total(params map[string]interface{}) int64 {
	db := database.GetDB()
	taskId, ok := params["TaskId"]
	if ok && taskId.(int) > 0 {
		db = db.Where("task_id = ?", taskId)
	}
	protocol, ok := params["Protocol"]
	if ok && protocol.(int) > 0 {
		db = db.Where("protocol = ?", protocol)
	}
	status, ok := params["Status"]
	if ok && status.(int) > -1 {
		db = db.Where("status = ?", status)
	}
	var count int64
	err := db.Model(&taskLog).Count(&count)
	if err.Error != nil {
		logrus.Error(err.Error)
		return -1
	}
	return count
}
