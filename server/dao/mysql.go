package dao

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDB() {
	dsn := "root:123456@tcp(localhost:3306)/go-chess?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	//migration
	db.AutoMigrate(&User{})

	Db = db
}

func GetUsername(userId int64) (string, error) {
	var user User
	if err := Db.Where("user_id = ?", userId).First(&user).Error; err != nil {
		return "", err
	}
	return user.Username, nil
}
