package loginLogController

import (
	"github.com/gin-gonic/gin"
	"github.com/wuchunfu/JobFlow/common"
	"github.com/wuchunfu/JobFlow/model/loginLogModel"
	"net/http"
)

func Index(ctx *gin.Context) {
	page, pageSize := common.ParseQueryParams(ctx)
	username := ctx.Query("username")
	loginLog := new(loginLogModel.LoginLog)
	dataList, count := loginLog.List(page, pageSize, username)
	ctx.JSON(http.StatusOK, gin.H{
		"code":     http.StatusOK,
		"data":     &dataList,
		"msg":      "获取数据成功！",
		"page":     page,
		"pageSize": pageSize,
		"total":    count,
	})
}
