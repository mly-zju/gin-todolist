package service

import (
	"fmt"
	"gin-todolist/library/cache"
	"gin-todolist/library/redis"
	"gin-todolist/model"
	"gin-todolist/model/dao"
)

func getItemListKey(tagId int) string {
	return fmt.Sprintf(model.REDIS_ITEM_LIST_KEY, tagId)
}

func getItemListKeyField(page, pageSize int) string {
	return fmt.Sprintf("%d_%d", page, pageSize)
}

// GetItemListByPage 分页获取item列表
func GetItemListByPage(tagId, page, pageSize int) ([]*dao.Item, error) {
	var result []*dao.Item
	// 先从redis获取
	key := getItemListKey(tagId)
	err := cache.HGetCacheData(key, getItemListKeyField(page, pageSize), 3600, dao.NewItemDao(), "GetItemListByPage", []interface{}{tagId, page, pageSize}, &result)
	return result, err
}

// EditItem 编辑item
func EditItem(id int, title, content, img string, tagIdList []int) error {
	_, err := dao.NewItemDao().EditItem(id, title, content, img, tagIdList)
	if err != nil {
		return err
	}
	// 删除tagId的缓存
	keys := make([]string, len(tagIdList))
	if len(keys) != 0 {
		for index, tagId := range tagIdList {
			keys[index] = getItemListKey(tagId)
		}
		rds := redis.GetRedis()
		err = rds.Del(keys...).Err()
	}
	return err
}
