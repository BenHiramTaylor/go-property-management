package main

import (
	"log"
	"time"

	"github.com/BenHiramTaylor/go-property-management/database"
	"github.com/BenHiramTaylor/go-property-management/properties"
	"github.com/BenHiramTaylor/go-property-management/tennants"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func initialiseDB() error {
	var err error
	database.DBConn, err = gorm.Open(sqlite.Open("production.db"), &gorm.Config{})
	if err != nil {
		return err
	}
	err = database.DBConn.Table("Properties").AutoMigrate(&properties.Property{})
	if err != nil {
		return err
	}
	err = database.DBConn.Table("Tennants").AutoMigrate(&tennants.Tennant{})
	if err != nil {
		return err
	}
	return nil
}

func initialiseRoutes() *fiber.App {
	app := fiber.New(fiber.Config{ReadTimeout: 60 * time.Second, WriteTimeout: 120 * time.Second, IdleTimeout: 12 * time.Hour})
	app.Get("/properties", properties.GetAllProperties)
	app.Post("/properties", properties.AddProperty)
	app.Get("/properties/:id", properties.GetIndividualProperty)
	app.Delete("/properties/:id", properties.DeleteProperty)
	return app
}

func main() {
	err := initialiseDB()
	if err != nil {
		panic(err)
	}
	log.Println("Initialised DB")
	app := initialiseRoutes()
	log.Fatalln(app.Listen(":80"))
}
