package cache

import (
	"encoding/json"
	"gin-todolist/library"
	"gin-todolist/library/redis"
	"reflect"
	"time"
)

func getCacheData(key string, field string, ttl time.Duration, daoObj interface{}, methodName string, params []interface{}, target interface{}) error {
	// 首先尝试从redis获取数据
	rds := redis.GetRedis()
	// 根据field是否为空来判断是否hash
	isHash := field != ""
	var bytes []byte
	var err error
	if isHash {
		bytes, err = rds.HGet(key, field).Bytes()
	} else {
		bytes, err = rds.Get(key).Bytes()
	}

	if rds.IsErr(err) {
		// 如果读取报错，返回报错
		return err
	} else if !rds.IsNil(err) {
		// 如果结果不为空，反序列化到target并返回
		library.WriteNotice("cache hit with key: " + key)
		json.Unmarshal(bytes, target)
		return nil
	}

	// 如果为空，需要从数据库读取
	methodV := reflect.ValueOf(daoObj).MethodByName(methodName)
	// 将params转为reflect.Value
	paramsV := formatReflectParams(params)
	// 执行数据库读取方法
	resultV := methodV.Call(paramsV)
	// 查看返回结果第二个参数是否nil
	if err, ok := resultV[1].Interface().(error); ok {
		return err
	}
	// 如果读取没有错误，将结果写回redis
	resultData := resultV[0].Interface()
	resBytes, err := json.Marshal(resultData)
	if isHash {
		err = rds.HSet(key, field, resBytes).Err()
		if int64(ttl) > 0 {
			rds.Expire(key, ttl)
		}
	} else {
		err = rds.Set(key, resBytes, ttl).Err()
	}
	// 将结果赋值给target
	json.Unmarshal(resBytes, target)
	return err
}

func GetCacheData(key string, ttl int, daoObj interface{}, methodName string, params []interface{}, target interface{}) error {
	iTTL := time.Duration(ttl) * time.Second
	return getCacheData(key, "", iTTL, daoObj, methodName, params, target)
}

func HGetCacheData(key string, field string, ttl int, daoObj interface{}, methodName string, params []interface{}, target interface{}) error {
	iTTL := time.Duration(ttl) * time.Second
	return getCacheData(key, field, iTTL, daoObj, methodName, params, target)
}

func formatReflectParams(params []interface{}) []reflect.Value {
	// 将params转为reflect.Value
	paramsV := []reflect.Value{}
	for _, param := range params {
		paramsV = append(paramsV, reflect.ValueOf(param))
	}
	return paramsV
}
