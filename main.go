package main

import (
	"gin-vue/config"
	"gin-vue/middleware/database"
	"gin-vue/routers"
	"gin-vue/service/taskService"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	initConfig := config.InitConfig
	db := database.InitDB()
	defer db.Close()

	taskService.ServiceTask.Initialize()

	router := routers.InitRouter()

	port := initConfig.System.Host
	if port != "" {
		panic(router.Run(port))
	}
	// listen and serve on 0.0.0.0:8080
	panic(router.Run())
}
