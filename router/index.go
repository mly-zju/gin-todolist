package router

import (
	"gin-todolist/controller"
	"gin-todolist/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter(server *gin.Engine) {
	// 初始化路由和中间件
	// 统计时间中间件
	server.Use(middleware.Timer())

	server.GET("/taglist", controller.TagList.HandleFunc)
	server.GET("/itemlist", controller.ItemList.HandleFunc)
	server.POST("/itemedit", controller.ItemEdit.HandleFunc)
	server.Static("/static", "./upload")
}
