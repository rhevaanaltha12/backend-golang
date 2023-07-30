package config

import (
	"backend-golang/helper"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {

	var err error

	env := godotenv.Load()

	if env != nil {
		// helpers.Logger("error", "Error Getting Env")
		helper.Logger("error", "Error getting env")
	}

	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, dbHost, dbPort, dbName)

	db, err = gorm.Open(mysql.Open(dbURI), &gorm.Config{})

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
	}

	sql, err := db.DB()

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return
	}

	sql.SetMaxOpenConns(100)
	sql.SetMaxIdleConns(10)
	sql.SetConnMaxLifetime(1800000)
	sql.SetConnMaxIdleTime(600000)

	fmt.Println("Success | connected db MySQL")
	fmt.Println(db)

	// db.Debug().AutoMigrate(model.User{})
	// db.Debug().AutoMigrate(model.Seller{})
	// db.Debug().AutoMigrate(model.Fish{})
	// db.Debug().AutoMigrate(model.Temporary{})

}

func GetDB() *gorm.DB {
	return db
}
