package base

import (
	"fmt"
	"gin-todolist/library"
	"sync"
)

type Errcode struct {
	Name        string
	Errno       int
	ErrMsg      string
	InnerErrMsg string
}

var (
	mu         sync.Mutex
	errCodeMap map[string]*Errcode
)

// Errcode实现了Error接口
func (e *Errcode) Error() string {
	return fmt.Sprintf("api level error: %s, inner error: %s", e.ErrMsg, e.InnerErrMsg)
}

// GetError 根据配置文件中的错误名生成错误对象
func GetError(errName string) *Errcode {
	if errCodeMap == nil {
		mu.Lock()
		defer mu.Unlock()
		if errCodeMap == nil {
			// 从配置文件读取数据
			type errConf struct {
				List []*Errcode
			}
			var errconf errConf
			library.GetConf("errcode", &errconf)
			// 将数组转为map
			errCodeMap = make(map[string]*Errcode)
			for _, item := range errconf.List {
				errCodeMap[item.Name] = item
			}
		}
	}

	res, ok := errCodeMap[errName]
	if !ok {
		return &Errcode{
			Errno:       500,
			InnerErrMsg: "Internal error",
		}
	}
	return res
}
