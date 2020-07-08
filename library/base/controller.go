package base

import (
	"errors"
	"gin-todolist/library"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Implement 子controller需要实现的接口
type Implement interface {
	Validate(c *gin.Context) bool
	Execute(c *gin.Context) (interface{}, error)
}

// BaseController controller基类，模板方法的实现
type BaseController struct {
	Implement
}

func NewBaseController(imp Implement) *BaseController {
	return &BaseController{
		Implement: imp,
	}
}

func (bctl *BaseController) EjsonErr(c *gin.Context, err error) {
	switch err.(type) {
	case *Errcode:
		errcode := err.(*Errcode)
		c.JSON(http.StatusOK, gin.H{
			"errno":  errcode.Errno,
			"errmsg": errcode.ErrMsg,
		})
	default:
		// 打印错误
		c.JSON(http.StatusOK, gin.H{
			"errno":  500,
			"errmsg": "Internal error",
		})
	}
	// 打印错误
	library.WriteWarn(err)
}

// Handle gin要求的模板方法，子类只需要填充接口实现即可
func (bctl *BaseController) HandleFunc(c *gin.Context) {
	// 首先recover，即所有panic都在这里挡一下
	defer func() {
		if r := recover(); r != nil {
			// err可能有很多类型：string，error，其他
			var err error
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Internal error")
			}

			bctl.EjsonErr(c, err)
		}
	}()

	if !bctl.Validate(c) {
		bctl.EjsonErr(c, GetError("paramsError"))
		return
	}

	data, err := bctl.Execute(c)
	if err != nil {
		bctl.EjsonErr(c, err)
		return
	}
	// 否则，最后返回正常结果
	c.JSON(http.StatusOK, gin.H{
		"errno": 0,
		"data":  data,
	})
}
