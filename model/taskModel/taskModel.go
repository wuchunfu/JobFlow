package taskModel

import (
	"gin-vue/common"
	"gin-vue/middleware/database"
	"gin-vue/model/taskHostModel"
	logger "github.com/sirupsen/logrus"
	"strings"
	//"github.com/robfig/cron/v3"
)

type TaskProtocol int8

const (
	TaskHTTP int8 = iota + 1 // HTTP协议
	TaskRPC                  // RPC方式执行命令
)

type TaskLevel int8

const (
	TaskLevelParent int8 = 1 // 父任务
	TaskLevelChild  int8 = 2 // 子任务(依赖任务)
)

type TaskDependencyStatus int8

const (
	TaskDependencyStatusStrong int8 = 1 // 强依赖
	TaskDependencyStatusWeak   int8 = 2 // 弱依赖
)

type TaskHTTPMethod int8

const (
	TaskHTTPMethodGet  int8 = 1
	TaskHttpMethodPost int8 = 2
)

// 任务
type Task struct {
	TaskId           int                      `json:"taskId" gorm:"type:int(20); primary_key; auto_increment; not null"`
	TaskName         string                   `json:"taskName" gorm:"type:varchar(100); not null"`                              // 任务名称
	TaskLevel        int8                     `json:"taskLevel" gorm:"type:int(11); not null; index:IDX_task_level; default 1"` // 任务等级 1: 主任务 2: 依赖任务
	DependencyTaskId string                   `json:"dependencyTaskId" gorm:"type:varchar(100); not null; default ''"`          // 依赖任务ID,多个ID逗号分隔
	DependencyStatus int8                     `json:"dependencyStatus" gorm:"type:int(11); not null; default 1"`                // 依赖关系 1:强依赖 主任务执行成功, 依赖任务才会被执行 2:弱依赖
	CronExpression   string                   `json:"cronExpression" gorm:"type:varchar(100); not null"`                        // crontab
	Protocol         int8                     `json:"protocol" gorm:"type:int(11); not null; index:IDX_task_protocol"`          // 协议 1:http 2:系统命令
	Command          string                   `json:"command" gorm:"type:varchar(300); not null"`                               // URL地址或shell命令
	HttpMethod       int8                     `json:"httpMethod" gorm:"type:int(11); not null; default 1"`                      // http请求方法
	Timeout          int                      `json:"timeout" gorm:"type:int(11); not null; default 0"`                         // 任务执行超时时间(单位秒),0不限制
	IsMultiInstance  int8                     `json:"isMultiInstance" gorm:"type:int(11) ; not null; default 1"`                // 是否允许多实例运行
	RetryTimes       int8                     `json:"retryTimes" gorm:"type:int(11); not null; default 0"`                      // 重试次数
	RetryInterval    int16                    `json:"retryInterval" gorm:"type:int(11); not null; default 0"`                   // 重试间隔时间
	NotifyStatus     int8                     `json:"notifyStatus" gorm:"type:int(11); not null; default 1"`                    // 任务执行结束是否通知 0: 不通知 1: 失败通知 2: 执行结束通知 3: 任务执行结果关键字匹配通知
	NotifyType       int8                     `json:"notifyType" gorm:"type:int(11); not null; default 0"`                      // 通知类型 1: 邮件 2: slack 3: webhook
	NotifyReceiverId string                   `json:"notifyReceiverId" gorm:"type:varchar(300); not null; default ''"`          // 通知接受者ID, setting表主键ID，多个ID逗号分隔
	NotifyKeyword    string                   `json:"notifyKeyword" gorm:"type:varchar(200); not null; default ''"`
	TaskTag          string                   `json:"taskTag" gorm:"type:varchar(100); not null; default ''"`
	TaskRemark       string                   `json:"taskRemark" gorm:"type:varchar(200); not null; default ''"`                   // 备注
	TaskStatus       int8                     `json:"taskStatus" gorm:"type:int(11); not null; index:IDX_task_status ; default 0"` // 状态 1:正常 0:停止
	CreateTime       string                   `json:"createTime" gorm:"type:varchar(50); not null"`                                // 创建时间
	UpdateTime       string                   `json:"updateTime" gorm:"type:varchar(50); DEFAULT ''"`                              // 更新时间
	DeleteTime       string                   `json:"deleteTime" gorm:"type:varchar(50); DEFAULT ''"`                              // 删除时间
	Hosts            []taskHostModel.TaskHost `json:"hosts" gorm:"-"`
	NextRunTime      string                   `json:"nextRunTime" gorm:"-"`
}

func (task *Task) List(page int, pageSize int, taskName string) ([]Task, int) {
	db := database.GetDB()
	if taskName != "" {
		db = db.Where("task_name = ?", taskName)
	}
	list := make([]Task, 0)
	findErr := db.Model(&task).Offset((page - 1) * pageSize).Limit(pageSize).Order("update_time desc").Find(&list)
	if findErr.Error != nil {
		logger.Error(findErr.Error)
		return nil, -1
	}
	var count int
	countErr := db.Model(&task).Count(&count)
	if countErr.Error != nil {
		logger.Error(countErr.Error)
		return nil, -1
	}
	return list, count
}

func (task *Task) Detail(taskId int) *Task {
	db := database.GetDB()
	err := db.Model(&task).Where("task_id = ?", taskId).Find(&task)
	if err.Error != nil {
		logger.Error(err.Error)
		return nil
	}
	taskHost := new(taskHostModel.TaskHost)
	task.Hosts = taskHost.GetHostIdsByTaskId(taskId)
	return task
}

// 新增
func (task *Task) Save() int {
	db := database.GetDB().Begin()
	err := db.Model(&task).Create(&task)
	if err.Error != nil {
		db.Rollback()
		logger.Error(err.Error)
	}
	db.Commit()
	return task.TaskId
}

// 删除
func (task *Task) Delete(taskId int) {
	db := database.GetDB().Begin()
	err := db.Model(&task).Where("task_id = ?", taskId).Delete(&task)
	if err.Error != nil {
		db.Rollback()
		logger.Error(err.Error)
	}
	db.Commit()
}

func (task *Task) IsExistsTaskName(taskName string) int {
	db := database.GetDB()
	var count int
	err := db.Model(&task).Where("task_name = ?", taskName).Count(&count)
	if err.Error != nil {
		logger.Error(err.Error)
		return -1
	}
	return count
}
func (task *Task) GetStatus(taskId int) int8 {
	db := database.GetDB()
	err := db.Model(&task).Where("task_id = ?", taskId).Find(&task)
	if err.Error != nil {
		logger.Error(err.Error)
		return 0
	}
	return task.TaskStatus
}

// 获取某个主机下的所有激活任务
func (task *Task) ActiveListByHostId(hostId int) ([]Task, error) {
	db := database.GetDB()
	taskHost := new(taskHostModel.TaskHost)
	taskIds, err := taskHost.GetTaskIdsByHostId(hostId)
	if err != nil {
		return nil, err
	}
	if len(taskIds) == 0 {
		return nil, nil
	}
	list := make([]Task, 0)
	aa := make([]interface{}, 0)
	aa = append(aa, taskIds...)
	bb := make([]int64, 0)
	for _, value := range bb {
		bb = append(bb, value)
	}
	findErr := db.Model(&task).Where("task_status = ? AND task_level = ? AND id IN (?)", common.Enabled, TaskLevelParent, bb).Find(&list)
	if findErr.Error != nil {
		return list, findErr.Error
	}
	return task.setHostsForTasks(list)
}

func (task *Task) setHostsForTasks(tasks []Task) ([]Task, error) {
	taskHost := new(taskHostModel.TaskHost)
	var err error
	for i, value := range tasks {
		taskHostDetails := taskHost.GetHostIdsByTaskId(value.TaskId)
		tasks[i].Hosts = taskHostDetails
	}
	return tasks, err
}

// 激活
func (task *Task) Enable(taskId int) int {
	db := database.GetDB().Begin()
	fieldMap := make(map[string]interface{})
	fieldMap["task_status"] = common.Enabled
	err := db.Model(&task).Where("task_id = ?", taskId).Update(fieldMap)
	if err.Error != nil {
		db.Rollback()
		logger.Error(err.Error)
	}
	db.Commit()
	return task.TaskId
}

// 停止
func (task *Task) Disable(taskId int) int {
	db := database.GetDB().Begin()
	fieldMap := make(map[string]interface{})
	fieldMap["task_status"] = common.Disabled
	err := db.Model(&task).Where("task_id = ?", taskId).Update(fieldMap)
	if err.Error != nil {
		db.Rollback()
		logger.Error(err.Error)
	}
	db.Commit()
	return task.TaskId
}

// 获取所有激活任务
func (task *Task) ActiveList(page int, pageSize int) ([]Task, error) {
	db := database.GetDB()
	list := make([]Task, 0)
	err := db.Model(&task).Where("task_status = ? AND task_level = ?", common.Enabled, TaskLevelParent).Offset((page - 1) * pageSize).Limit(pageSize).Find(&list)
	if err.Error != nil {
		return list, err.Error
	}
	return task.setHostsForTasks(list)
}

// 获取依赖任务列表
func (task *Task) GetDependencyTaskList(ids string) ([]Task, error) {
	list := make([]Task, 0)
	if ids == "" {
		return list, nil
	}
	idList := strings.Split(ids, ",")
	taskIds := make([]interface{}, len(idList))
	for i, v := range idList {
		taskIds[i] = v
	}
	db := database.GetDB()
	err := db.Raw("select t.* from gin_task where t.level = ? and t.id in (?)").Find(&list)
	if err.Error != nil {
		return list, err.Error
	}
	return task.setHostsForTasks(list)
}
