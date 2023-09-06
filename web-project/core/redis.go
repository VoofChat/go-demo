package core

import (
	"gorm-web/conf"
	"gorm-web/pkg/redis"
)

var RedisCoreClient *redis.Redis

func InitRedis() {
	for key, c := range conf.BasicConf.Redis {
		var err error
		switch key {
		case "dxcore":
			{
				RedisCoreClient, err = redis.InitRedisClient(c)
				break
			}
		}
		if err != nil {
			panic("init redis failed!")
		}
	}
}

func CloseRedis() {
	if nil != RedisCoreClient {
		_ = RedisCoreClient.Close()
	}
}
