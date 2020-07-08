package controller

import (
	"gin-todolist/library"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// @todo 这里修改为配置文件
const UPLOAD_DIR = "./upload"

// ImgUploadHandler 上传文件
func ImgUploadHandler(c *gin.Context) {
	file, _ := c.FormFile("picture")
	library.WriteNotice("get picture with name: " + file.Filename)

	// 计算文件存储路径
	tmpFileName := strconv.Itoa(int(time.Now().Unix())) + path.Ext(file.Filename)
	dst := UPLOAD_DIR + "/" + tmpFileName
	err := c.SaveUploadedFile(file, dst)
	if err != nil {
		library.CheckErr(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"errno": true,
		"url":   "localhost:8080/static/" + tmpFileName,
	})
}
