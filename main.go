package main

import (
	"fmt"
	"graduate/acam"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := "root:123456@tcp(127.0.0.1:3306)/word_val?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}
	fmt.Println(db)
	match := acam.NewMatcher()
	match.ReBuild(db)
}
