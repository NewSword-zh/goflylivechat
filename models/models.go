package models

import (
	"fmt"
	"goflylivechat/common"
	"log"
	"time"

	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

// 数据库基础字段
type Model struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

func init() {
	Connect()
}
func Connect() error {
	mysql := common.GetMysqlConf()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysql.Username, mysql.Password, mysql.Server, mysql.Port, mysql.Database)
	var err error
	DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		log.Println(err)
		panic("数据库连接失败!")
		return err
	}
	DB.SingularTable(true)                       // 禁用表名复数化
	DB.LogMode(true)                             // 在控制台打印执行的 SQL 语句、执行时间、影响行数等日志
	DB.DB().SetMaxIdleConns(10)                  // 最大空闲连接数
	DB.DB().SetMaxOpenConns(100)                 // 最大打开连接数
	DB.DB().SetConnMaxLifetime(59 * time.Second) // 连接最大存活时间 MySQL 的 wait_timeout 默认为 60 秒，这里设为 59 秒可提前关闭，避免连接失效
	return nil
}
func Execute(sql string) error {
	return DB.Exec(sql).Error
}
func CloseDB() {
	defer DB.Close()
}
