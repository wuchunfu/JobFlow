package hostController

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wuchunfu/JobFlow/common"
	"github.com/wuchunfu/JobFlow/middleware/rpc/client"
	"github.com/wuchunfu/JobFlow/middleware/rpc/grpcpool"
	rpc "github.com/wuchunfu/JobFlow/middleware/rpc/proto"
	"github.com/wuchunfu/JobFlow/model/hostModel"
	"github.com/wuchunfu/JobFlow/model/taskModel"
	"github.com/wuchunfu/JobFlow/service/taskService"
	"github.com/wuchunfu/JobFlow/utils/datetimeUtils"
	"net/http"
	"strconv"
	"strings"
)

type HostForm struct {
	HostId     int
	HostIds    []int
	HostAlias  string
	HostName   string
	HostPort   int
	Remark     string
	CreateTime string
	UpdateTime string
}

// Index 用户列表页
func Index(ctx *gin.Context) {
	page, pageSize := common.ParseQueryParams(ctx)
	hostName := ctx.Query("hostName")
	host := new(hostModel.Host)
	dataList, count := host.List(page, pageSize, hostName)
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
	//hostId, _ := strconv.ParseInt(ctx.Param("hostId"), 10, 64)
	hostId, _ := strconv.Atoi(ctx.Param("hostId"))
	host := new(hostModel.Host)
	detail := host.Detail(hostId)
	if detail.HostId >= 0 {
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

func Save(ctx *gin.Context) {
	form := new(HostForm)
	ctx.Bind(form)
	hostId := form.HostId
	hostAlias := strings.TrimSpace(form.HostAlias)
	hostName := strings.TrimSpace(form.HostName)
	hostPort := form.HostPort
	remark := strings.TrimSpace(form.Remark)

	host := new(hostModel.Host)
	hostAliasCount := host.IsExistsHostAlias(hostAlias)
	hostNameCount := host.IsExistsHostName(hostName)
	if hostAliasCount > 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"data": nil,
			"msg":  "主机别名已存在！",
		})
	} else if hostNameCount > 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"data": nil,
			"msg":  "主机名已存在！",
		})
	} else {
		host.HostAlias = hostAlias
		host.HostName = hostName
		host.HostPort = hostPort
		host.Remark = remark
		host.CreateTime = datetimeUtils.FormatDateTime()
		host.UpdateTime = datetimeUtils.FormatDateTime()
		host.Save()

		addr := fmt.Sprintf("%s:%d", hostName, hostPort)
		grpcpool.Pool.Release(addr)

		task := new(taskModel.Task)
		tasks, err := task.ActiveListByHostId(hostId)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code": http.StatusOK,
				"data": nil,
				"msg":  "刷新任务主机信息失败！",
			})
		}
		taskService.ServiceTask.BatchAdd(tasks)

		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"data": nil,
			"msg":  "保存成功！",
		})
	}
}

func Update(ctx *gin.Context) {
	form := new(HostForm)
	ctx.Bind(form)
	hostId := form.HostId
	hostAlias := strings.TrimSpace(form.HostAlias)
	hostName := strings.TrimSpace(form.HostName)
	hostPort := form.HostPort
	remark := strings.TrimSpace(form.Remark)

	host := new(hostModel.Host)
	hostMap := make(map[string]interface{})
	hostMap["host_alias"] = hostAlias
	hostMap["host_name"] = hostName
	hostMap["host_port"] = hostPort
	hostMap["remark"] = remark
	hostMap["update_time"] = datetimeUtils.FormatDateTime()
	host.Update(hostId, hostMap)
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": nil,
		"msg":  "修改成功！",
	})
}

func Delete(ctx *gin.Context) {
	form := new(HostForm)
	ctx.Bind(form)
	host := new(hostModel.Host)
	for _, id := range form.HostIds {
		host.Delete(id)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": nil,
		"msg":  "删除成功！",
	})
}

func Ping(ctx *gin.Context) {
	hostId, _ := strconv.Atoi(ctx.Param("hostId"))

	host := new(hostModel.Host)
	hostInfo := host.Detail(hostId)
	if hostInfo == nil || host.HostId <= 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"data": nil,
			"msg":  "主机不存在",
		})
	}

	const testConnectionCommand = "echo hello"
	const testConnectionTimeout = 5

	taskReq := new(rpc.TaskRequest)
	taskReq.Command = testConnectionCommand
	taskReq.Timeout = testConnectionTimeout
	output, err := client.Exec(host.HostName, host.HostPort, taskReq)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"data": nil,
			"msg":  "连接失败-" + output + err.Error(),
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"data": nil,
			"msg":  "连接成功",
		})
	}
}

//设置默认路由当访问一个错误网站时返回
func NotFound(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusNotFound,
		"data": nil,
		"msg":  "404 ,page not exists!",
	})
}
