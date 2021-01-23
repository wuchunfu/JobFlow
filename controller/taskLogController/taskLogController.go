package taskLogController

import (
	"github.com/gin-gonic/gin"
	"github.com/wuchunfu/JobFlow/common"
	"github.com/wuchunfu/JobFlow/model/taskLogModel"
	"net/http"
)

// Index 用户列表页
func Index(ctx *gin.Context) {
	page, pageSize := common.ParseQueryParams(ctx)
	taskName := ctx.Query("taskName")
	taskLog := new(taskLogModel.TaskLog)
	dataList, count := taskLog.List(page, pageSize, taskName)
	ctx.JSON(http.StatusOK, gin.H{
		"code":     http.StatusOK,
		"data":     &dataList,
		"msg":      "获取数据成功！",
		"page":     page,
		"pageSize": pageSize,
		"total":    count,
	})
}
