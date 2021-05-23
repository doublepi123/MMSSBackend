package util

import (
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	DB    *gorm.DB
	Redis *redis.Client
}

//初始化数据库连接并检查数据库连接
//若此处检查连接失败，则会直接panic
func (db *Database) Init() {
	var err error

	// import "gorm.io/driver/mysql"
	// refer: https://gorm.io/docs/connecting_to_the_database.html#MySQL
	dsn := os.Getenv("DSN")
	db.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println(err)
		return
	}
	sqlDB, err := db.DB.DB()
	if err != nil {
		fmt.Println(err)
		return
	}
	// 使用连接池
	sqlDB.SetMaxIdleConns(200)
	sqlDB.SetMaxOpenConns(200)
	sqlDB.SetConnMaxIdleTime(time.Hour)
	sqlDB.SetConnMaxLifetime(time.Hour)
	db.Redis = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       2,  // use DB 2
	})

	err = db.Redis.Set("test", "test", 0).Err()
	if err != nil {
		panic(err)
	}
}
