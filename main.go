package main

import (
	"fmt"
	"log"

	"github.com/BenHiramTaylor/go-property-management/properties"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func GetDatabase() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("production.db"), &gorm.Config{})
	if err != nil {
		return &gorm.DB{}, fmt.Errorf("failed to connect database")
	}
	return db, nil
}

func initialiseDB() error {
	db, err := GetDatabase()
	if err != nil {
		return err
	}
	err = db.Table("Properties").AutoMigrate(&properties.Property{})
	if err != nil {
		return err
	}
	return nil
}

func main() {
	log.Println("Starting Server...")
	err := initialiseDB()
	if err != nil {
		panic(err)
	}
}
