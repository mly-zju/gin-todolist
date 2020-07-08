package redis

import (
	"fmt"
	"gin-todolist/library"

	"github.com/go-redis/redis"
)

type redisConfig struct {
	IP       string
	Port     int
	Password string
}

const Nil = redis.Nil

type Redis struct {
	*redis.Client
}

// GetRedis: 获取redis连接
func GetRedis() *Redis {
	var rdsconf redisConfig

	if _, err := library.GetConf("redis", &rdsconf); err != nil {
		library.CheckErr(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", rdsconf.IP, rdsconf.Port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &Redis{
		client,
	}
}

// IsNil: 是否为空
func (rds *Redis) IsNil(err error) bool {
	return err == Nil
}

// IsErr: 是否为非空错误
func (rds *Redis) IsErr(err error) bool {
	return err != nil && err != Nil
}
