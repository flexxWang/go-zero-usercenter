package svc

import (
	"database/sql"
	"time"

	"gozero-user-demo/internal/config"
	"gozero-user-demo/internal/middleware"
	"gozero-user-demo/internal/model"

	_ "github.com/go-sql-driver/mysql"
	redisstore "github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config    config.Config
	DB        *sql.DB
	Redis     *redisstore.Redis
	UserModel *model.UserModel
	Audit     rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := mustOpenDB(c)
	rdb := redisstore.MustNewRedis(c.Redis)
	userModel := model.NewUserModel(db, rdb, c.Cache.UserProfileSeconds)

	return &ServiceContext{
		Config:    c,
		DB:        db,
		Redis:     rdb,
		UserModel: userModel,
		Audit:     middleware.NewAuditMiddleware(rdb, c.Session.TokenPrefix).Handle,
	}
}

func mustOpenDB(c config.Config) *sql.DB {
	db, err := sql.Open("mysql", c.MySQL.DataSource)
	if err != nil {
		panic(err)
	}

	if c.MySQL.MaxOpenConns > 0 {
		db.SetMaxOpenConns(c.MySQL.MaxOpenConns)
	}
	if c.MySQL.MaxIdleConns > 0 {
		db.SetMaxIdleConns(c.MySQL.MaxIdleConns)
	}
	if c.MySQL.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(time.Duration(c.MySQL.ConnMaxLifetime) * time.Minute)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	return db
}
