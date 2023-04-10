package configs

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func Connect() (*gorm.DB, error) {
	database, err := gorm.Open(mysql.Open("root:@tcp(localhost:3306)/katchum_tecnology"), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MySQL!")
	db = database

	return db, nil
}
