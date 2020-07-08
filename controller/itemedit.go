package controller

import (
	"gin-todolist/library/base"
	"gin-todolist/model/service"

	"github.com/gin-gonic/gin"
)

type ItemEditParams struct {
	Id        int    `form:"id" binding:"required"`
	Title     string `form:"title" binding:"required`
	Content   string `form:"content" binding:"required"`
	Img       string `form:"img"`
	TagIdList []int  `form:"tagIdList" binding:"required"`
}

type ItemEditController struct {
	*base.BaseController
	ItemEditParams
}

func (ctl *ItemEditController) Validate(c *gin.Context) bool {
	var reqParams ItemEditParams
	if err := c.ShouldBindJSON(&reqParams); err != nil {
		return false
	}
	ctl.ItemEditParams = reqParams
	return true
}

func (ctl *ItemEditController) Execute(c *gin.Context) (interface{}, error) {
	err := service.EditItem(ctl.ItemEditParams.Id, ctl.ItemEditParams.Title, ctl.ItemEditParams.Content, ctl.ItemEditParams.Img, ctl.ItemEditParams.TagIdList)
	return true, err
}

func NewItemEditController() *ItemEditController {
	ItemEdit := &ItemEditController{}
	ItemEdit.BaseController = base.NewBaseController(ItemEdit)
	return ItemEdit
}

var ItemEdit = NewItemEditController()
