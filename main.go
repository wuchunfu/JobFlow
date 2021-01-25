package main

import (
	"github.com/wuchunfu/JobFlow/config"
	"github.com/wuchunfu/JobFlow/middleware/database"
	"github.com/wuchunfu/JobFlow/routers"
	"github.com/wuchunfu/JobFlow/service/taskService"
)

func main() {
	initConfig := config.InitConfig
	database.InitDB()

	taskService.ServiceTask.Initialize()

	router := routers.InitRouter()

	port := initConfig.System.Host
	if port != "" {
		panic(router.Run(port))
	}
	// listen and serve on 0.0.0.0:8080
	panic(router.Run())
}
