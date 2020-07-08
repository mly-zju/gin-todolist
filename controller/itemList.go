package controller

import (
	"gin-todolist/library/base"
	"gin-todolist/model/service"

	"github.com/gin-gonic/gin"
)

type RequestParams struct {
	Id       int `form:"id" binding:"required"`
	Page     int `form:"page"`
	PageSize int `form:"pageSize"`
}

type ItemListController struct {
	*base.BaseController
	RequestParams
}

func (ctl *ItemListController) Validate(c *gin.Context) bool {
	var reqParams RequestParams
	if err := c.ShouldBindQuery(&reqParams); err != nil {
		return false
	}

	// 检测page和pageSize
	if reqParams.Page == 0 {
		reqParams.Page = 1
	}
	if reqParams.PageSize == 0 {
		reqParams.PageSize = 10
	}
	ctl.RequestParams = reqParams
	return true
}

func (ctl *ItemListController) Execute(c *gin.Context) (interface{}, error) {
	return service.GetItemListByPage(ctl.Id, ctl.Page, ctl.PageSize)
}

func NewItemListController() *ItemListController {
	ItemList := &ItemListController{}
	ItemList.BaseController = base.NewBaseController(ItemList)
	return ItemList
}

var ItemList = NewItemListController()
