module scrapy

go 1.17

require (
	gorm.io/driver/mysql v1.3.3
	gorm.io/gorm v1.23.5
	platform v0.0.0-00010101000000-000000000000
)

require (
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.4 // indirect
)

replace platform => ../platform
