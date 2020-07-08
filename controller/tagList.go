package controller

import (
	"gin-todolist/library/base"
	"gin-todolist/model/service"

	"github.com/gin-gonic/gin"
)

type TagListController struct {
	*base.BaseController
}

func (ctl *TagListController) Validate(c *gin.Context) bool {
	return true
}

func (ctl *TagListController) Execute(c *gin.Context) (interface{}, error) {
	return service.GetTagList()
}

func NewTagListController() *TagListController {
	TagList := &TagListController{}
	TagList.BaseController = base.NewBaseController(TagList)
	return TagList
}

var TagList = NewTagListController()
