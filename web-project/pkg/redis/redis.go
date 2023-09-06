package redis

import (
	env2 "gorm-web/pkg/env"
	log2 "gorm-web/pkg/log"
	"gorm-web/utils"
	"time"

	"github.com/gin-gonic/gin"
	redigo "github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
)

// 日志打印Do args部分支持的最大长度
const logForRedisValue = 50

type RedisConf struct {
	Service         string        `yaml:"service"`
	Addr            string        `yaml:"addr"`
	Password        string        `yaml:"password"`
	MaxIdle         int           `yaml:"maxIdle"`
	MaxActive       int           `yaml:"maxActive"`
	IdleTimeout     time.Duration `yaml:"idleTimeout"`
	MaxConnLifetime time.Duration `yaml:"maxConnLifetime"`
	ConnTimeOut     time.Duration `yaml:"connTimeOut"`
	ReadTimeOut     time.Duration `yaml:"readTimeOut"`
	WriteTimeOut    time.Duration `yaml:"writeTimeOut"`
}

func (conf *RedisConf) checkConf() {
	if conf.MaxIdle == 0 {
		conf.MaxIdle = 50
	}
	if conf.MaxActive == 0 {
		conf.MaxActive = 100
	}
	if conf.IdleTimeout == 0 {
		conf.IdleTimeout = 5 * time.Minute
	}
	if conf.MaxConnLifetime == 0 {
		conf.MaxConnLifetime = 10 * time.Minute
	}
	if conf.ConnTimeOut == 0 {
		conf.ConnTimeOut = 3 * time.Second
	}
	if conf.ReadTimeOut == 0 {
		conf.ReadTimeOut = 1200 * time.Millisecond
	}
	if conf.WriteTimeOut == 0 {
		conf.WriteTimeOut = 1200 * time.Millisecond
	}
}

// Redis 日志打印Do args部分支持的最大长度
type Redis struct {
	pool       *redigo.Pool
	service    string
	remoteAddr string
	logger     *zap.Logger
}

func InitRedisClient(conf RedisConf) (*Redis, error) {
	conf.checkConf()

	p := &redigo.Pool{
		MaxIdle:         conf.MaxIdle,
		MaxActive:       conf.MaxActive,
		IdleTimeout:     conf.IdleTimeout,
		MaxConnLifetime: conf.MaxConnLifetime,
		Wait:            true,
		Dial: func() (conn redigo.Conn, e error) {
			con, err := redigo.Dial(
				"tcp",
				conf.Addr,
				redigo.DialPassword(conf.Password),
				redigo.DialConnectTimeout(conf.ConnTimeOut),
				redigo.DialReadTimeout(conf.ReadTimeOut),
				redigo.DialWriteTimeout(conf.WriteTimeOut),
			)
			if err != nil {
				return nil, err
			}
			return con, nil
		},
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
	c := &Redis{
		service:    conf.Service,
		remoteAddr: conf.Addr,
		pool:       p,
		logger:     log2.GetLogger().WithOptions(zap.AddCallerSkip(1)),
	}

	return c, nil
}

func (r *Redis) Do(ctx *gin.Context, commandName string, args ...interface{}) (reply interface{}, err error) {
	start := time.Now()

	conn := r.pool.Get()
	if err := conn.Err(); err != nil {
		r.logger.Error("get connection error: "+err.Error(), r.commonFields(ctx)...)
		return reply, err
	}

	reply, err = conn.Do(commandName, args...)
	if e := conn.Close(); e != nil {
		r.logger.Warn("connection close error: "+e.Error(), r.commonFields(ctx)...)
	}

	end := time.Now()

	// 执行时间 单位:毫秒
	ralCode := 0
	msg := "redis do success"
	if err != nil {
		ralCode = -1
		msg = "redis do error: " + err.Error()
		r.logger.Error(msg)
	}

	fields := append(r.commonFields(ctx),
		zap.String("requestStartTime", utils.GetFormatRequestTime(start)),
		zap.String("requestEndTime", utils.GetFormatRequestTime(end)),
		zap.Float64("cost", utils.GetRequestCost(start, end)),
		zap.String("command", commandName),
		zap.String("commandVal", utils.JoinArgs(logForRedisValue, args...)),
		zap.Int("ralCode", ralCode),
	)

	r.logger.Info(msg, fields...)
	return reply, err
}

func (r *Redis) commonFields(ctx *gin.Context) []zap.Field {
	return []zap.Field{
		zap.String("prot", "redis"),
		zap.String("module", env2.GetAppName()),
		zap.String("localIp", env2.LocalIP),
		zap.String("remoteAddr", r.remoteAddr),
		zap.String("service", r.service),
		zap.String("uri", ctx.GetString(env2.ContextKeyUri)),
	}
}

func (r *Redis) Close() error {
	return r.pool.Close()
}

func (r *Redis) Stats() (inUseCount, idleCount, activeCount int) {
	stats := r.pool.Stats()
	idleCount = stats.IdleCount
	activeCount = stats.ActiveCount
	inUseCount = activeCount - idleCount
	return inUseCount, idleCount, activeCount
}
