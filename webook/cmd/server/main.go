package main

import (
	"Project-WeBook/webook/config"
	"Project-WeBook/webook/internal/handler"
	"Project-WeBook/webook/internal/repository"
	localCache "Project-WeBook/webook/internal/repository/cache"
	"Project-WeBook/webook/internal/repository/dao"
	"Project-WeBook/webook/internal/router"
	"Project-WeBook/webook/internal/service"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		// panic 相当于整个 goroutine 结束
		panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}

	cache := localCache.NewUserCache(rdb)
	ud := dao.NewUserDAO(db)
	repo := repository.NewUserRepository(ud, cache)
	svc := service.NewUserService(repo)
	u := handler.NewUserHandler(svc)
	server := router.RegisterRoutes(u)

	err = server.Run(":8082")
	if err != nil {
		return
	}
}
