package common

import (
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"strconv"
)

type CommonMap map[string]interface{}

const (
	IsAdmin  int8 = 0 // 普通用户
	Disabled int8 = 0 // 禁用
	Failure  int8 = 0 // 失败
	Enabled  int8 = 1 // 启用
	Running  int8 = 1 // 运行中
	Finish   int8 = 2 // 完成
	Cancel   int8 = 3 // 取消
)

const (
	Page        = 1    // 当前页数
	PageSize    = 10   // 每页多少条数据
	MaxPageSize = 1000 // 每次最多取多少条
)

// 解析查询参数中的页数和每页数量
func ParseQueryParams(ctx *gin.Context) (int, int) {
	//page, pageError := strconv.ParseInt(ctx.DefaultQuery("page", strconv.Itoa(Page)), 10, 64)
	page, pageError := strconv.Atoi(ctx.DefaultQuery("page", strconv.Itoa(Page)))
	if pageError != nil {
		logger.Error(pageError.Error())
	}
	//pageSize, pageSizeError := strconv.ParseInt(ctx.DefaultQuery("pageSize", strconv.Itoa(PageSize)), 10, 64)
	pageSize, pageSizeError := strconv.Atoi(ctx.DefaultQuery("pageSize", strconv.Itoa(PageSize)))
	if pageSizeError != nil {
		logger.Error(pageSizeError.Error())
	}
	return page, pageSize
}

type Response struct {
	Code     int16       `json:"code"`
	Data     interface{} `json:"data"`
	Msg      string      `json:"msg"`
	Page     int64       `json:"page"`
	PageSize int64       `json:"pageSize"`
	Total    int64       `json:"total"`
}

type BaseModel struct {
	Page     int64 `json:"page" gorm:"-"`
	PageSize int64 `json:"pageSize" gorm:"-"`
}

func (baseModel *BaseModel) ParsePageAndPageSize(params CommonMap) {
	page, ok := params["Page"]
	if ok {
		baseModel.Page = page.(int64)
	}
	pageSize, ok := params["PageSize"]
	if ok {
		baseModel.PageSize = pageSize.(int64)
	}
	if baseModel.Page <= 0 {
		baseModel.Page = Page
	}
	if baseModel.PageSize <= 0 {
		baseModel.PageSize = MaxPageSize
	}
}

func (baseModel *BaseModel) PageLimitOffset() int64 {
	return (baseModel.Page - 1) * baseModel.PageSize
}
