package db

import (
	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"message-service/internal/pkg/config"
)

var Db *gorm.DB

func init() {
	// https://gorm.io/docs/connecting_to_the_database.html#Existing-database-connection
	mysqlConfig := config.GetMySqlConfig()
	// username:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	url := mysqlConfig.User + ":" + mysqlConfig.Password + "@tcp(" + mysqlConfig.Host + ":" + mysqlConfig.Port + ")/" + mysqlConfig.Dbname + "?charset=utf8mb4&parseTime=True&loc=Local"
	conn, err := sql.Open("mysql", url)
	if err != nil {
		panic("failed to connect database")
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: conn,
	}), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	Db = db
}
