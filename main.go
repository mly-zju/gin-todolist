package main

import (
	"gin-todolist/router"
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// 日志
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	server := gin.Default()
	router.InitRouter(server)

	server.Run(":8080")
}
