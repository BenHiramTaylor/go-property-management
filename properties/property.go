package properties

import (
	"gorm.io/gorm"
)

type Property struct {
	gorm.Model
	PropertyType        string
	Address             string
	NumberOfBedrooms    uint
	PurchasePriceDollar uint
}
