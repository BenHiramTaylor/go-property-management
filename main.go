package go_property_management

import (
	"fmt"

	"github.com/BenHiramTaylor/go-property-management/properties"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	fmt.Println("Starting Server...")
	db, err := gorm.Open(sqlite.Open("properties.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&properties.Property{})

	// Create
	db.Create(&properties.Property{PropertyType: "Apartment", Address: "Steward Building, Steward Street", NumberOfBedrooms: 2, PurchasePriceDollar: 250000})
}
