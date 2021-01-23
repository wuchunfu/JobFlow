package loginLogController

import (
	"gin-vue/common"
	"gin-vue/model/loginLogModel"
	"github.com/gin-gonic/gin"
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
