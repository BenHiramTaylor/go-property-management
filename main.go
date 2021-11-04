package main

import (
	"log"
	"time"

	"github.com/BenHiramTaylor/go-property-management/database"
	"github.com/BenHiramTaylor/go-property-management/properties"
	"github.com/BenHiramTaylor/go-property-management/tennants"
	"github.com/BenHiramTaylor/go-property-management/users"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// INITIALISE DATABASE AND TABLES
func initialiseDB() error {
	var err error
	var tables = map[string]interface{}{
		"Properties": &properties.Property{},
		"Tennants":   &tennants.Tennant{},
		"Users":      &users.User{},
	}
	database.DBConn, err = gorm.Open(sqlite.Open("production.db"), &gorm.Config{})
	if err != nil {
		return err
	}
	for table, schema := range tables {
		err = database.DBConn.Table(table).AutoMigrate(schema)
		if err != nil {
			return err
		}
	}
	err = users.CreateDefaultAdmin("sys_admin", "MrR0b0t123$")
	if err != nil {
		return err
	}
	return nil
}

func initialiseRoutes() *fiber.App {
	// CREATE APP WITH BASE CONFIG
	app := fiber.New(fiber.Config{ReadTimeout: 60 * time.Second, WriteTimeout: 120 * time.Second, IdleTimeout: 12 * time.Hour})

	// ADD AUTHENTICATION MIDDLEWARE
	app.Use(basicauth.New(basicauth.Config{Authorizer: users.CheckAuth}))

	// REGISTER PROPERTY ENDPOINTS
	app.Get("/properties", properties.GetAllProperties)
	app.Post("/properties", properties.AddProperty)
	app.Get("/properties/:id", properties.GetIndividualProperty)
	app.Put("/properties/:id", properties.UpdateProperty)
	app.Delete("/properties/:id", properties.DeleteProperty)

	// REGISTER TENNANT ENDPOINTS
	app.Get("/tennants", tennants.GetAllTennants)
	app.Post("/tennants", tennants.AddTennant)
	app.Get("/tennants/:id", tennants.GetIndividualTennant)
	app.Put("/tennants/:id", tennants.UpdateTennant)
	app.Delete("/tennants/:id", tennants.DeleteTennant)
	app.Post("/tennants/:tennantID/properties/:propertyID", tennants.AssignTennantToProperty)

	// REGISTER USER ENDPOINTS
	app.Post("/users", users.AddUser)
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
