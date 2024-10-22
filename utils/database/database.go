package database

import (
	"fmt"
	"skripsi/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(c config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DB_User,
		c.DB_Password,
		c.DB_Host,
		c.DB_Port,
		c.DB_Name,
	)
	fmt.Println("DSN:", dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
