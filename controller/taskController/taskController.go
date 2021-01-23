package taskController

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/wuchunfu/JobFlow/common"
	"github.com/wuchunfu/JobFlow/model/hostModel"
	"github.com/wuchunfu/JobFlow/model/taskHostModel"
	"github.com/wuchunfu/JobFlow/model/taskModel"
	"github.com/wuchunfu/JobFlow/service/taskService"
	"github.com/wuchunfu/JobFlow/utils"
	"github.com/wuchunfu/JobFlow/utils/datetimeUtils"
	"net/http"
	"strconv"
	"strings"
)

// 任务
type TaskForm struct {
	TaskId           int
	TaskIds          []int
	TaskName         string
	TaskLevel        int8
	DependencyTaskId string
	DependencyStatus int8
	CronExpression   string
	Protocol         int8
	Command          string
	HttpMethod       int8
	HostId           string
	Timeout          int
	IsMultiInstance  int8
	RetryTimes       int8
	RetryInterval    int16
	NotifyStatus     int8
	NotifyType       int8
	NotifyReceiverId string
	NotifyKeyword    string
	TaskTag          string
	TaskRemark       string
	TaskStatus       int8
	CreateTime       string
	UpdateTime       string
	DeleteTime       string
	Hosts            []taskHostModel.TaskHost
	NextRunTime      string
}

// Index 用户列表页
func Index(ctx *gin.Context) {
	page, pageSize := common.ParseQueryParams(ctx)
	taskName := ctx.Query("taskName")
	task := new(taskModel.Task)
	dataList, count := task.List(page, pageSize, taskName)
	for i, item := range dataList {
		nextRunTime := taskService.ServiceTask.NextRunTime(item)
		dataList[i].NextRunTime = nextRunTime
		//dataList[i].NextRunTime = taskService.ServiceTask.NextRunTime(item)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":     http.StatusOK,
		"data":     &dataList,
		"msg":      "获取数据成功！",
		"page":     page,
		"pageSize": pageSize,
		"total":    count,
	})
}

// Detail 用户详情
func Detail(ctx *gin.Context) {
	//taskId, _ := strconv.ParseInt(ctx.Param("taskId"), 10, 64)
	taskId, _ := strconv.Atoi(ctx.Param("taskId"))
	task := new(taskModel.Task)
	detail := task.Detail(taskId)
	if detail.TaskId >= 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"data": &detail,
			"msg":  "获取数据成功！",
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"data": &detail,
			"msg":  "获取数据失败！",
		})
	}
}

// Index 用户列表页
func HostAllList(ctx *gin.Context) {
	host := new(hostModel.Host)
	dataList := host.AllList()
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": &dataList,
		"msg":  "获取数据成功！",
	})
}

func Save(ctx *gin.Context) {
	form := new(TaskForm)
	ctx.Bind(form)
	taskId := form.TaskId
	taskName := strings.TrimSpace(form.TaskName)
	taskLevel := form.TaskLevel
	dependencyTaskId := strings.TrimSpace(form.DependencyTaskId)
	dependencyStatus := form.DependencyStatus
	cronExpression := strings.TrimSpace(form.CronExpression)
	protocol := form.Protocol
	httpMethod := form.HttpMethod
	command := strings.TrimSpace(form.Command)
	timeout := form.Timeout
	isMultiInstance := form.IsMultiInstance
	retryTimes := form.RetryTimes
	retryInterval := form.RetryInterval
	notifyStatus := form.NotifyStatus
	notifyType := form.NotifyType
	notifyReceiverId := strings.TrimSpace(form.NotifyReceiverId)
	notifyKeyword := strings.TrimSpace(form.NotifyKeyword)
	taskTag := strings.TrimSpace(form.TaskTag)
	taskRemark := strings.TrimSpace(form.TaskRemark)

	task := new(taskModel.Task)
	taskNameCount := task.IsExistsTaskName(taskName)
	if taskNameCount > 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"data": nil,
			"msg":  "任务名已存在！",
		})
	} else if form.Protocol == taskModel.TaskRPC && form.HostId == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"data": nil,
			"msg":  "请选择主机名！",
		})
	} else {
		task.TaskName = taskName
		task.TaskLevel = taskLevel
		task.DependencyTaskId = dependencyTaskId
		task.DependencyStatus = dependencyStatus
		task.CronExpression = cronExpression
		task.Protocol = protocol
		task.HttpMethod = httpMethod
		task.Command = command
		task.Timeout = timeout
		task.IsMultiInstance = isMultiInstance
		task.RetryTimes = retryTimes
		task.RetryInterval = retryInterval
		task.NotifyStatus = notifyStatus - 1
		task.NotifyType = notifyType - 1
		task.NotifyReceiverId = notifyReceiverId
		task.NotifyKeyword = notifyKeyword
		task.TaskTag = taskTag
		task.TaskRemark = taskRemark
		task.TaskStatus = common.Running
		task.CreateTime = datetimeUtils.FormatDateTime()
		task.UpdateTime = datetimeUtils.FormatDateTime()

		if task.IsMultiInstance != 1 {
			task.IsMultiInstance = 0
		}
		if task.NotifyStatus > 0 && task.NotifyType != 3 && task.NotifyReceiverId == "" {
			ctx.JSON(http.StatusOK, gin.H{
				"code": http.StatusBadRequest,
				"data": nil,
				"msg":  "至少选择一个通知接收者！",
			})
		}
		if task.Protocol == taskModel.TaskHTTP {
			command := strings.ToLower(task.Command)
			if !strings.HasPrefix(command, "http://") && !strings.HasPrefix(command, "https://") {
				ctx.JSON(http.StatusOK, gin.H{
					"code": http.StatusBadRequest,
					"data": nil,
					"msg":  "请输入正确的URL地址！",
				})
			}
			if task.Timeout > 300 {
				ctx.JSON(http.StatusOK, gin.H{
					"code": http.StatusBadRequest,
					"data": nil,
					"msg":  "HTTP任务超时时间不能超过300秒！",
				})
			}
		}
		if task.RetryTimes > 10 || task.RetryTimes < 0 {
			ctx.JSON(http.StatusOK, gin.H{
				"code": http.StatusBadRequest,
				"data": nil,
				"msg":  "任务重试次数取值0-10！",
			})
		}
		if task.RetryInterval > 3600 || task.RetryInterval < 0 {
			ctx.JSON(http.StatusOK, gin.H{
				"code": http.StatusBadRequest,
				"data": nil,
				"msg":  "任务重试间隔时间取值0-3600！",
			})
		}
		if task.DependencyStatus != taskModel.TaskDependencyStatusStrong &&
			task.DependencyStatus != taskModel.TaskDependencyStatusWeak {
			ctx.JSON(http.StatusOK, gin.H{
				"code": http.StatusBadRequest,
				"data": nil,
				"msg":  "请选择依赖关系！",
			})
		}

		if task.TaskLevel == taskModel.TaskLevelParent {
			err := utils.PanicToError(func() {
				cron.ParseStandard(form.CronExpression)
			})
			if err != nil {
				ctx.JSON(http.StatusOK, gin.H{
					"code": http.StatusBadRequest,
					"data": nil,
					"msg":  "crontab表达式解析失败！",
				})
			}
		} else {
			task.DependencyTaskId = ""
			task.CronExpression = ""
		}
		if taskId > 0 && task.DependencyTaskId != "" {
			dependencyTaskIds := strings.Split(task.DependencyTaskId, ",")
			if utils.InStringSlice(dependencyTaskIds, strconv.Itoa(taskId)) {
				ctx.JSON(http.StatusOK, gin.H{
					"code": http.StatusBadRequest,
					"data": nil,
					"msg":  "不允许设置当前任务为子任务！",
				})
			}
		}

		taskId := task.Save()

		taskHost := new(taskHostModel.TaskHost)
		if form.Protocol == taskModel.TaskRPC {
			hostIdStrList := strings.Split(form.HostId, ",")
			hostIds := make([]int, len(hostIdStrList))
			for i, hostIdStr := range hostIdStrList {
				hostIds[i], _ = strconv.Atoi(hostIdStr)
			}
			taskHost.Save(taskId, hostIds)
		} else {
			taskHost.Delete(taskId)
		}

		status := task.GetStatus(taskId)
		if status == common.Enabled && task.TaskLevel == taskModel.TaskLevelParent {
			addTaskToTimer(taskId)
		}

		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"data": nil,
			"msg":  "保存成功！",
		})
	}
}

func Delete(ctx *gin.Context) {
	form := new(TaskForm)
	ctx.Bind(form)
	taskHost := new(taskHostModel.TaskHost)
	task := new(taskModel.Task)
	for _, id := range form.TaskIds {
		task.Delete(id)
		taskHost.Delete(id)
		taskService.ServiceTask.Remove(id)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": nil,
		"msg":  "删除成功！",
	})
}

// 添加任务到定时器
func addTaskToTimer(taskId int) {
	task := new(taskModel.Task)
	detail := task.Detail(taskId)
	taskService.ServiceTask.RemoveAndAdd(*detail)
}

// 手动运行任务
func Run(ctx *gin.Context) {
	taskId, _ := strconv.Atoi(ctx.Param("taskId"))

	task := new(taskModel.Task)
	detail := task.Detail(taskId)

	//detail.CronExpression = "手动运行"
	taskService.ServiceTask.Run(*detail)

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": nil,
		"msg":  "任务已开始运行, 请到任务日志中查看结果！",
	})
}

// 停止运行中的任务
func Stop(ctx *gin.Context) {
	taskId, _ := strconv.Atoi(ctx.Param("taskId"))
	task := new(taskModel.Task)
	detail := task.Detail(taskId)
	if detail.Protocol != taskModel.TaskRPC {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"data": nil,
			"msg":  "仅支持SHELL任务手动停止！",
		})
	} else if len(detail.Hosts) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"data": nil,
			"msg":  "任务节点列表为空！",
		})
	} else {
		for _, host := range detail.Hosts {
			taskService.ServiceTask.Stop(host.HostName, host.HostPort, int64(taskId))
		}
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"data": nil,
			"msg":  "已执行停止操作, 请等待任务退出！",
		})
	}
}

// 激活任务
func Enable(ctx *gin.Context) {
	taskId, _ := strconv.Atoi(ctx.Param("taskId"))
	task := new(taskModel.Task)
	enable := task.Enable(taskId)

	if enable < 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"data": nil,
			"msg":  "操作失败！",
		})
	}

	addTaskToTimer(taskId)

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": nil,
		"msg":  "操作成功！",
	})
}

// 暂停任务
func Disable(ctx *gin.Context) {
	taskId, _ := strconv.Atoi(ctx.Param("taskId"))
	task := new(taskModel.Task)
	disable := task.Disable(taskId)

	if disable < 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"data": nil,
			"msg":  "操作失败！",
		})
	}

	taskService.ServiceTask.RemoveAndStop(taskId)

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": nil,
		"msg":  "操作成功！",
	})
}
