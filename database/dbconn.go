// https://dev.to/siddheshk02/jwt-authentication-in-go-5dp7
package database

import (
	"fmt"
	"log"

	"github.com/AlwiLion/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "admin123" //Enter your password for the DB
	dbname   = "bankdb"
)

var dsn string = fmt.Sprintf("host=%s port=%d user=%s "+
	"password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
	host, port, user, password, dbname)

var DB *gorm.DB

func DBconn() {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	DB = db

	db.AutoMigrate(&models.User{}) // we are going to create a models.go file for the User Model.
	db.AutoMigrate(&models.UserStatment{})
}
