package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	MySQL struct {
		DataSource      string
		MaxOpenConns    int
		MaxIdleConns    int
		ConnMaxLifetime int
	}
	Redis redis.RedisConf
	Cache struct {
		UserProfileSeconds int
	}
	Session struct {
		TokenPrefix string
	}
}
