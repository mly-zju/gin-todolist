package dao

import (
	"gin-todolist/library/db"
)

// Tag的字段格式
type Tag struct {
	Id   int    `json:"id" ddb:"id"`
	Name string `json:"name" ddb:"name"`
}

type TagDao struct {
	TableName string
}

func NewTagDao() *TagDao {
	return &TagDao{
		TableName: "todolist_tag",
	}
}

func (dao *TagDao) GetTagList() ([]*Tag, error) {
	var result []*Tag
	err := db.GetConn().Select(dao.TableName, map[string]interface{}{"status": 1}, []string{"id", "name"}, &result)

	return result, err
}
