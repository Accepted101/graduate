package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	. "platform"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//按行读取
func WriteWord(db *gorm.DB) {
	filePath := "./word_base/webdict/webdict_with_freq.txt"
	file, err := os.OpenFile(filePath, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("open file error! :", err)
		panic(err)
	}
	defer file.Close()
	stat, err := file.Stat()
	fmt.Printf("file size = %d\n", stat.Size)
	buf := bufio.NewReader(file)
	wordslice := make([]WordVal, 0)
	last := time.Now()
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				if len(wordslice) > 0 {
					err := db.Create(wordslice).Error
					if err != nil {
						fmt.Println("insert into DB error :", err)
						return
					}
				}
				fmt.Println("data successfully modify!")
				return
			} else {
				panic(err)
			}
		}
		line = strings.TrimSuffix(line, "\r\n")
		slice := strings.Split(line, " ")
		val, err := strconv.Atoi(slice[1])
		if err != nil {
			panic(err)
		}
		wordslice = append(wordslice, WordVal{WordName: slice[0], WordVal: uint(val)})
		if len(wordslice) > 1000 {
			err := db.Create(wordslice).Error
			if err != nil {
				fmt.Println("insert into DB error :", err)
				return
			}
			fmt.Printf("insert successful : %d word\n", len(wordslice))
			fmt.Println("cost time : ", time.Now().Sub(last))
			last = time.Now()
			wordslice = make([]WordVal, 0)
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
	fmt.Println(db)
}
