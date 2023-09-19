package main

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

func Connect(user string, password string, host string, database string) (*gorm.DB, error) {
	if db == nil {
		dsn := user + ":" + password + "@" + "tcp(" + host + ")/" + database
		log.Println("connecting to the database")
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	}
	return db, err
}
