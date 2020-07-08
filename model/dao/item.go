package dao

import (
	"fmt"
	"gin-todolist/library"
	"gin-todolist/library/base"
	"gin-todolist/library/db"
)

// Item dao以item表为主体，同时涵盖了relation表
type Item struct {
	Id        int    `json:"id" ddb:"id"`
	Title     string `json:"title" ddb:"title"`
	Content   string `json:"content" ddb:"content"`
	Img       string `json:"img" ddb:"img"`
	CreatedAt string `json:"createdAt" ddb:"created_at"`
	UpdatedAt string `json:"updatedAt" ddb:"updated_at"`
}

type Relation struct {
	Id     int `json:"id" ddb:"id"`
	ItemId int `json:"itemId" ddb:"item_id"`
	TagId  int `json:"tagId" ddb:"tag_id"`
}

type ItemDao struct {
	ItemTableName     string
	RelationTableName string
}

func NewItemDao() *ItemDao {
	return &ItemDao{
		ItemTableName:     "todolist_item",
		RelationTableName: "todolist_relation",
	}
}

func (dao *ItemDao) GetItemListByPage(tagId, page, pageSize int) ([]*Item, error) {
	var result []*Item
	m, n := (page-1)*pageSize, pageSize
	preSql := fmt.Sprintf("select a.* from %s as a inner join %s as b on b.tag_id = {{tagId}} and a.id = b.item_id limit {{m}}, {{n}}", dao.ItemTableName, dao.RelationTableName)
	err := db.GetConn().NamedQuery(preSql, map[string]interface{}{
		"tagId": tagId,
		"m":     m,
		"n":     n,
	}, &result)
	return result, err
}

func (dao *ItemDao) EditItem(id int, title, content, img string, tagIdList []int) (int, error) {
	// 首先查询文章是否存在，如果不存在，报错
	var result []*Item
	dbCtx := db.GetConn()
	err := dbCtx.Select(dao.ItemTableName, map[string]interface{}{"id": id}, []string{"id"}, &result)
	if err != nil {
		return 0, err
	}
	if result == nil {
		return 0, base.GetError("itemNotExist")
	}

	// 其次查询已存在的tagList，对比新传来的，查看哪些需要增加，哪些要删除
	relationList := []*Relation{}
	err = dbCtx.Select(dao.RelationTableName, map[string]interface{}{"item_id": id}, []string{"id", "item_id", "tag_id"}, &relationList)
	if err != nil {
		return 0, err
	}
	// 查看需要增加的
	var addTagList []int
	for _, tagId := range tagIdList {
		isExist := false
		for _, relationItem := range relationList {
			if relationItem.TagId == tagId {
				isExist = true
				break
			}
		}

		if !isExist {
			addTagList = append(addTagList, tagId)
		}
	}
	// 查看需要删除的
	var delTagList []int
	for _, relationItem := range relationList {
		isExist := false
		for _, tagId := range tagIdList {
			if relationItem.TagId == tagId {
				isExist = true
				break
			}
		}

		if !isExist {
			delTagList = append(delTagList, relationItem.TagId)
		}
	}
	library.WriteNotice(delTagList)

	// 开始使用事务
	tx, err := dbCtx.BeginTx()
	if err != nil {
		return 0, err
	}
	// 1. 首先更新item表
	rowsAff, err := dbCtx.UpdateTx(tx, dao.ItemTableName, db.Where{"id": id}, db.RowData{
		"title":   title,
		"content": content,
		"img":     img,
	})
	if err != nil {
		dbCtx.Rollback(tx)
		return 0, err
	}
	// 2. 其次新增relation表
	if len(addTagList) != 0 {
		newData := make(db.RowDataArr, len(addTagList))
		for index, tagId := range addTagList {
			newData[index] = db.RowData{
				"item_id": id,
				"tag_id":  tagId,
			}
		}
		_, err = dbCtx.InsertTx(tx, dao.RelationTableName, newData)
		if err != nil {
			dbCtx.Rollback(tx)
			return 0, err
		}
	}
	// 3. 删除relation表
	if len(delTagList) != 0 {
		err = dbCtx.DeleteTx(tx, dao.RelationTableName, db.Where{
			"tag_id in": delTagList,
			"item_id":   id,
		})
		if err != nil {
			dbCtx.Rollback(tx)
			return 0, err
		}
	}
	// 提交
	dbCtx.Commit(tx)

	return int(rowsAff), err
}
