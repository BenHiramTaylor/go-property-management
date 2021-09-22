package main

import (
	"fmt"

	"github.com/BenHiramTaylor/go-property-management/properties"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	fmt.Println("Starting Server...")
	db, err := gorm.Open(sqlite.Open("production.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	propertiesTable := db.Table("properties")

	// Migrate the schema
	propertiesTable.AutoMigrate(&properties.Property{})

	// Create
	propertiesTable.Create(&properties.Property{UUID: uuid.New(), PropertyType: "Apartment", Address: "Steward Building, Steward Street", NumberOfBedrooms: 2, PurchasePriceDollar: 250000})
}
