package main

import (
	"fmt"
	"graduate/acam"
	. "graduate/platform"
	"math"
	"time"
	"unicode/utf8"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func modify(db *gorm.DB) {
	var id int64 = 0
	sum := 0
	word := make([]WordVal, 0, 1000)

	//拉取新数据
	fmt.Println("start rebuild: ", time.Now())
	for {
		err := db.Model(&WordVal{}).Limit(1000).Find(&word, "word_id > ?", id).Error
		if err != nil {
			panic(err)
		}
		if len(word) != 0 {
			id = word[len(word)-1].WordId
		} else {
			break
		}
		sum += len(word)
		fmt.Println("get words : ", len(word))
		for _, val := range word {
			numer := utf8.RuneCountInString(val.WordName)
			val.WordVal = int64(math.Pow(float64(val.Tmp), float64(numer)))
			err := db.Where("word_id = ?", val.WordId).Save(&val).Error
			if err != nil {
				panic(err)
			}
		}

	}
}

func main() {

	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := "root:123456@tcp(127.0.0.1:3306)/word_val?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	match := acam.NewMatcher()
	match.ReBuild(db)
	for {
		var text string
		fmt.Scanf("%s", &text)
		fmt.Println(text)
		match.FuncDynamicProgrammingAndDivide(text)
	}
}
