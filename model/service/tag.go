package service

import (
	"gin-todolist/library/cache"
	"gin-todolist/model"
	"gin-todolist/model/dao"
)

// GetTagList 获取tag列表
func GetTagList() ([]*dao.Tag, error) {
	var result []*dao.Tag
	// // 先从redis获取
	key := model.REDIS_TAG_LIST_KEY
	err := cache.GetCacheData(key, 3600, dao.NewTagDao(), "GetTagList", []interface{}{}, &result)
	return result, err
}
