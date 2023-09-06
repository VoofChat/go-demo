package conf

import (
	env2 "gorm-web/pkg/env"
	log2 "gorm-web/pkg/log"
	mysql2 "gorm-web/pkg/mysql"
	redis2 "gorm-web/pkg/redis"
)

type ServerConfig struct {
	Port string `yaml:"port"`
}

// Config 对应config.yaml
type Config struct {
	Redis  map[string]redis2.RedisConf
	Mysql  map[string]mysql2.MysqlConf
	Log    log2.LogConfig
	Server ServerConfig
}

var (
	// BasicConf 配置文件对应的全局变量
	BasicConf Config
)

func InitConf() {
	// 加载通用基础配置（必须）
	env2.LoadConf("config.yaml", env2.SubConfMount, &BasicConf)
}
