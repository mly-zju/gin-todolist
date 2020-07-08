package library

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
)

// CheckErr: 检查异常
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

// WriteNotice: 打notice日志
func WriteNotice(info ...interface{}) {
	infoArr := []interface{}{"[info]"}
	infoArr = append(infoArr, info...)
	fmt.Fprintln(gin.DefaultWriter, infoArr...)
}

// WriteWarn: 打warnning日志
func WriteWarn(info ...interface{}) {
	infoArr := []interface{}{"[warnning]"}
	infoArr = append(infoArr, info...)
	fmt.Fprintln(gin.DefaultWriter, infoArr...)
}

// GetConf: 获取conf
func GetConf(fileName string, v interface{}) (toml.MetaData, error) {
	return toml.DecodeFile("./conf/"+fileName+".toml", v)
}
