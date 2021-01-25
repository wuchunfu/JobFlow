package main

import (
	"github.com/wuchunfu/JobFlow/cmd"
	"github.com/wuchunfu/JobFlow/middleware/config"
	"github.com/wuchunfu/JobFlow/middleware/database"
	"github.com/wuchunfu/JobFlow/middleware/logutil"
	"github.com/wuchunfu/JobFlow/routers"
	"github.com/wuchunfu/JobFlow/service/taskService"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 检查或设置命令行参数
	cmd.Execute()

	// 将日志写入文件或打印到控制台
	logutil.InitLog(&config.ServerSetting.Log)
	// 初始化数据库连接
	database.InitDB(&config.ServerSetting.Database)
	//初始化定时任务
	taskService.ServiceTask.Initialize()
	// 初始化路由
	router := routers.InitRouter()
	setting := &config.ServerSetting.System
	// 启动
	router.Run(setting.HttpAddr + ":" + setting.HttpPort)
}
