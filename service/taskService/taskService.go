package taskService

import (
	"fmt"
	"gin-vue/common"
	"gin-vue/config"
	"gin-vue/middleware/httpclient"
	rpcClient "gin-vue/middleware/rpc/client"
	rpc "gin-vue/middleware/rpc/proto"
	"gin-vue/model/taskHostModel"
	"gin-vue/model/taskLogModel"
	"gin-vue/model/taskModel"
	"gin-vue/utils"
	"gin-vue/utils/datetimeUtils"
	"github.com/robfig/cron/v3"
	//"github.com/jakecoffman/cron"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	ServiceTask Task
)

var (
	// 定时任务调度管理器
	serviceCron *cron.Cron
	//nyc, _ = time.LoadLocation("Asia/Shanghai")
	//serviceCron = cron.New(cron.WithSeconds(), cron.WithLocation(nyc))

	wg sync.WaitGroup

	// 同一任务是否有实例处于运行中
	runInstance Instance

	// 任务计数-正在运行的任务
	taskCount TaskCount

	// 并发队列, 限制同时运行的任务数量
	concurrencyQueue ConcurrencyQueue
)

// 并发队列
type ConcurrencyQueue struct {
	queue chan struct{}
}

func (cq *ConcurrencyQueue) Add() {
	cq.queue <- struct{}{}
}

func (cq *ConcurrencyQueue) Done() {
	<-cq.queue
}

// 任务计数
type TaskCount struct {
	wg   sync.WaitGroup
	exit chan struct{}
}

func (tc *TaskCount) Add() {
	tc.wg.Add(1)
}

func (tc *TaskCount) Done() {
	tc.wg.Done()
}

func (tc *TaskCount) Exit() {
	tc.wg.Done()
	<-tc.exit
}

func (tc *TaskCount) Wait() {
	tc.Add()
	tc.wg.Wait()
	close(tc.exit)
}

// 任务ID作为Key
type Instance struct {
	m sync.Map
}

// 是否有任务处于运行中
func (i *Instance) has(key int) bool {
	_, ok := i.m.Load(key)

	return ok
}

func (i *Instance) add(key int) {
	i.m.Store(key, struct{}{})
}

func (i *Instance) done(key int) {
	i.m.Delete(key)
}

type Task struct{}

type TaskResult struct {
	Result     string
	Err        error
	RetryTimes int8
}

// 初始化任务, 从数据库取出所有任务, 添加到定时任务并运行
func (task Task) Initialize() {
	setting := config.InitConfig
	//serviceCron = cron.New()
	nyc, _ := time.LoadLocation("Asia/Shanghai")
	serviceCron = cron.New(cron.WithSeconds(), cron.WithLocation(nyc))
	serviceCron.Start()
	concurrencyQueue = ConcurrencyQueue{queue: make(chan struct{}, setting.System.ConcurrencyQueue)}
	taskCount = TaskCount{sync.WaitGroup{}, make(chan struct{})}
	go taskCount.Wait()

	logrus.Info("开始初始化定时任务")
	newTask := new(taskModel.Task)
	taskNum := 0
	page := 1
	pageSize := 1000
	maxPage := 1000
	for page < maxPage {
		taskList, err := newTask.ActiveList(page, pageSize)
		if err != nil {
			logrus.Fatalf("定时任务初始化#获取任务列表错误: %s", err)
		}
		if len(taskList) == 0 {
			break
		}
		for _, item := range taskList {
			task.Add(item)
			taskNum++
		}
		page++
	}
	logrus.Infof("定时任务初始化完成, 共%d个定时任务添加到调度器", taskNum)
}

// 批量添加任务
func (task Task) BatchAdd(tasks []taskModel.Task) {
	for _, item := range tasks {
		task.RemoveAndAdd(item)
	}
}

// 删除任务后添加
func (task Task) RemoveAndAdd(taskModel taskModel.Task) {
	task.Remove(taskModel.TaskId)
	task.Add(taskModel)
}

// 删除任务后添加
func (task Task) RemoveAndStop(taskId int) {
	task.Remove(taskId)
	task.StopJob()
}

// 添加任务
func (task Task) Add(taskModels taskModel.Task) {
	if taskModels.TaskLevel == taskModel.TaskLevelChild {
		logrus.Errorf("添加任务失败#不允许添加子任务到调度器#任务Id-%d", taskModels.TaskId)
		return
	}
	//wg.Add(1)
	taskFunc := createJob(taskModels)
	if taskFunc == nil {
		logrus.Error("创建任务处理Job失败,不支持的任务协议#", taskModels.Protocol)
		return
	}
	//cronName := strconv.Itoa(taskModels.TaskId)
	err := utils.PanicToError(func() {
		//serviceCron.AddFunc(taskModels.CronExpression, taskFunc, cronName)
		serviceCron.AddFunc(taskModels.CronExpression, taskFunc)
		//serviceCron.AddFunc("*/2 * * * * *", func() {
		//	logrus.Println("schedule every two seconds ...")
		//})
	})
	if err != nil {
		logrus.Error("添加任务到调度器失败#", err)
	}
	//serviceCron.Start()
	// 使用sync.WaitGroup保证主 goroutine 不退出。因为c.Start()中只是启动了一个 goroutine，如果主 goroutine 退出了，整个程序就停止了。
	//wg.Wait()
	// 阻塞主线程不退出
	//select{}
}

func (task Task) NextRunTime(taskModels taskModel.Task) string {
	if taskModels.TaskLevel != taskModel.TaskLevelParent ||
		taskModels.TaskStatus != common.Enabled {
		return ""
	}
	entries := serviceCron.Entries()

	//taskName := strconv.Itoa(taskModels.TaskId)
	//for _, item := range entries {
	//	if item.Name == taskName {
	//		return item.Next
	//	}
	//}
	for _, item := range entries {
		return item.Next.Format("2006-01-02 15:04:05")
	}
	return ""
}

// 停止运行中的任务
func (task Task) Stop(ip string, port int, id int64) {
	rpcClient.Stop(ip, port, id)
}

func (task Task) Remove(taskId int) {
	serviceCron.Remove(cron.EntryID(taskId))
	//serviceCron.RemoveJob(strconv.Itoa(id))
}

func (task Task) StopJob() {
	serviceCron.Stop()
}

// 等待所有任务结束后退出
func (task Task) WaitAndExit() {
	serviceCron.Stop()
	taskCount.Exit()
}

// 直接运行任务
func (task Task) Run(taskModel taskModel.Task) {
	go createJob(taskModel)()
}

type Handler interface {
	Run(taskModel taskModel.Task, taskUniqueId int64) (string, error)
}

// HTTP任务
type HTTPHandler struct{}

// http任务执行时间不超过300秒
const HttpExecTimeout = 300

func (h *HTTPHandler) Run(task taskModel.Task, taskUniqueId int64) (result string, err error) {
	if task.Timeout <= 0 || task.Timeout > HttpExecTimeout {
		task.Timeout = HttpExecTimeout
	}
	var resp httpclient.ResponseWrapper
	if task.HttpMethod == taskModel.TaskHTTPMethodGet {
		resp = httpclient.Get(task.Command, task.Timeout)
	} else {
		urlFields := strings.Split(task.Command, "?")
		task.Command = urlFields[0]
		var params string
		if len(urlFields) >= 2 {
			params = urlFields[1]
		}
		resp = httpclient.PostParams(task.Command, params, task.Timeout)
	}
	// 返回状态码非200，均为失败
	if resp.StatusCode != http.StatusOK {
		return resp.Body, fmt.Errorf("HTTP状态码非200-->%d", resp.StatusCode)
	}

	return resp.Body, err
}

// RPC调用执行任务
type RPCHandler struct{}

func (h *RPCHandler) Run(taskModel taskModel.Task, taskUniqueId int64) (result string, err error) {
	taskRequest := new(rpc.TaskRequest)
	taskRequest.Timeout = int32(taskModel.Timeout)
	taskRequest.Command = taskModel.Command
	taskRequest.Id = taskUniqueId
	resultChan := make(chan TaskResult, len(taskModel.Hosts))
	for _, taskHost := range taskModel.Hosts {
		go func(th taskHostModel.TaskHost) {
			output, err := rpcClient.Exec(th.HostName, th.HostPort, taskRequest)
			errorMessage := ""
			if err != nil {
				errorMessage = err.Error()
			}
			outputMessage := fmt.Sprintf("主机: [%s-%s:%d]\n%s\n%s\n\n",
				th.HostAlias, th.HostName, th.HostPort, errorMessage, output,
			)
			resultChan <- TaskResult{Err: err, Result: outputMessage}
		}(taskHost)
	}

	var aggregationErr error = nil
	aggregationResult := ""
	for i := 0; i < len(taskModel.Hosts); i++ {
		taskResult := <-resultChan
		aggregationResult += taskResult.Result
		if taskResult.Err != nil {
			aggregationErr = taskResult.Err
		}
	}

	return aggregationResult, aggregationErr
}

// 创建任务日志
func createTaskLog(task taskModel.Task, status int8) int64 {
	taskLog := new(taskLogModel.TaskLog)
	taskLog.TaskId = task.TaskId
	taskLog.TaskName = task.TaskName
	taskLog.CronExpression = task.CronExpression
	taskLog.Protocol = task.Protocol
	taskLog.Command = task.Command
	taskLog.Timeout = task.Timeout
	if task.Protocol == taskModel.TaskRPC {
		aggregationHost := ""
		for _, host := range task.Hosts {
			aggregationHost += fmt.Sprintf("%s - %s<br>", host.HostAlias, host.HostName)
		}
		taskLog.HostName = aggregationHost
	}
	//parseStartTime, _ := time.Parse("2006-01-02 15:04:05", datetimeUtils.FormatDateTime())
	//currentDateTime, _ := time.Parse("2006-01-02 15:04:05", time.Now().String())
	taskLog.StartTime = datetimeUtils.FormatDateTime()
	taskLog.EndTime = datetimeUtils.FormatDateTime()
	taskLog.Status = status
	insertId := taskLog.Create()
	return insertId
}

// 更新任务日志
func updateTaskLog(taskLogId int64, taskResult TaskResult) int64 {
	taskLog := new(taskLogModel.TaskLog)
	var status int8
	result := taskResult.Result
	if taskResult.Err != nil {
		status = common.Failure
	} else {
		status = common.Finish
	}
	return taskLog.Update(taskLogId, common.CommonMap{
		"retry_times": taskResult.RetryTimes,
		"status":      status,
		"result":      result,
	})
}

func createJob(task taskModel.Task) cron.FuncJob {
	handler := createHandler(task)
	if handler == nil {
		return nil
	}
	taskFunc := func() {
		taskCount.Add()
		defer taskCount.Done()

		taskLogId := beforeExecJob(task)
		if taskLogId <= 0 {
			return
		}

		if task.IsMultiInstance == 0 {
			runInstance.add(task.TaskId)
			defer runInstance.done(task.TaskId)
		}

		concurrencyQueue.Add()
		defer concurrencyQueue.Done()

		logrus.Infof("开始执行任务#%s#命令-%s", task.TaskName, task.Command)
		taskResult := execJob(handler, task, taskLogId)
		logrus.Infof("任务完成#%s#命令-%s", task.TaskName, task.Command)
		afterExecJob(task, taskResult, taskLogId)
	}
	return taskFunc
}

func createHandler(task taskModel.Task) Handler {
	var handler Handler = nil
	switch task.Protocol {
	case taskModel.TaskHTTP:
		handler = new(HTTPHandler)
	case taskModel.TaskRPC:
		handler = new(RPCHandler)
	}
	return handler
}

// 任务前置操作
func beforeExecJob(task taskModel.Task) int64 {
	if task.IsMultiInstance == 0 && runInstance.has(task.TaskId) {
		return createTaskLog(task, common.Cancel)
	}
	taskLogId := createTaskLog(task, common.Running)
	logrus.Debugf("任务命令-%s", task.Command)
	return taskLogId
}

// 任务执行后置操作
func afterExecJob(task taskModel.Task, taskResult TaskResult, taskLogId int64) {
	updateTaskLog(taskLogId, taskResult)
	// 发送邮件
	go SendNotification(task, taskResult)
	// 执行依赖任务
	go execDependencyTask(task, taskResult)
}

// 执行依赖任务, 多个任务并发执行
func execDependencyTask(task taskModel.Task, taskResult TaskResult) {
	// 父任务才能执行子任务
	if task.TaskLevel != taskModel.TaskLevelParent {
		return
	}

	// 是否存在子任务
	dependencyTaskId := strings.TrimSpace(task.DependencyTaskId)
	if dependencyTaskId == "" {
		return
	}

	// 父子任务关系为强依赖, 父任务执行失败, 不执行依赖任务
	if task.DependencyStatus == taskModel.TaskDependencyStatusStrong && taskResult.Err != nil {
		logrus.Infof("父子任务为强依赖关系, 父任务执行失败, 不运行依赖任务#主任务ID-%d", task.TaskId)
		return
	}

	// 获取子任务
	model := new(taskModel.Task)
	tasks, err := model.GetDependencyTaskList(dependencyTaskId)
	if err != nil {
		logrus.Errorf("获取依赖任务失败#主任务ID-%d#%s", task.TaskId, err.Error())
		return
	}
	if len(tasks) == 0 {
		logrus.Errorf("依赖任务列表为空#主任务ID-%d", task.TaskId)
	}
	for _, task := range tasks {
		task.CronExpression = fmt.Sprintf("依赖任务(主任务ID-%d)", task.TaskId)
		ServiceTask.Run(task)
	}
}

// 发送任务结果通知
func SendNotification(task taskModel.Task, taskResult TaskResult) {
	var statusName string
	// 未开启通知
	if task.NotifyStatus == 0 {
		return
	}
	if task.NotifyStatus == 3 {
		// 关键字匹配通知
		if !strings.Contains(taskResult.Result, task.NotifyKeyword) {
			return
		}
	}
	if task.NotifyStatus == 1 && taskResult.Err == nil {
		// 执行失败才发送通知
		return
	}
	if task.NotifyType != 3 && task.NotifyReceiverId == "" {
		return
	}
	if taskResult.Err != nil {
		statusName = "失败"
	} else {
		statusName = "成功"
	}
	fmt.Println("=====" + statusName)
	//// 发送通知
	//msg := notify.Message{
	//	"task_type":        task.NotifyType,
	//	"task_receiver_id": task.NotifyReceiverId,
	//	"name":             task.TaskName,
	//	"output":           taskResult.Result,
	//	"status":           statusName,
	//	"task_id":          task.TaskId,
	//}
	//notify.Push(msg)
}

// 执行具体任务
func execJob(handler Handler, task taskModel.Task, taskUniqueId int64) TaskResult {
	defer func() {
		if err := recover(); err != nil {
			logrus.Error("panic#service/taskService.go:execJob#", err)
		}
	}()
	// 默认只运行任务一次
	var execTimes int8 = 1
	if task.RetryTimes > 0 {
		execTimes += task.RetryTimes
	}
	var i int8 = 0
	var output string
	var err error
	for i < execTimes {
		output, err = handler.Run(task, taskUniqueId)
		if err == nil {
			return TaskResult{Result: output, Err: err, RetryTimes: i}
		}
		i++
		if i < execTimes {
			logrus.Warnf("任务执行失败#任务id-%d#重试第%d次#输出-%s#错误-%s", task.TaskId, i, output, err.Error())
			if task.RetryInterval > 0 {
				time.Sleep(time.Duration(task.RetryInterval) * time.Second)
			} else {
				// 默认重试间隔时间，每次递增1分钟
				time.Sleep(time.Duration(i) * time.Minute)
			}
		}
	}
	return TaskResult{Result: output, Err: err, RetryTimes: task.RetryTimes}
}
