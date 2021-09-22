package properties

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Property struct {
	gorm.Model
	UUID                uuid.UUID
	PropertyType        string
	Address             string
	NumberOfBedrooms    uint
	PurchasePriceDollar uint
}
