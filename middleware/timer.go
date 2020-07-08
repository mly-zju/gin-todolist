package middleware

import (
	"fmt"
	"gin-todolist/library"
	"time"

	"github.com/gin-gonic/gin"
)

func Timer() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		// 计算执行耗时
		processTime := time.Since(startTime).Microseconds()
		library.WriteNotice(fmt.Sprintf("spend time: %d", processTime))
	}
}
