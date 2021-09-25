package properties

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Property struct {
	gorm.Model
	ID                  uuid.UUID
	PropertyType        string `json:"property_type"`
	Address             string `json:"address"`
	NumberOfBedrooms    uint   `json:"number_of_bedrooms"`
	PurchasePriceDollar uint   `json:"purchase_price_dollar"`
}

func GetDatabase() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("production.db"), &gorm.Config{})
	if err != nil {
		return &gorm.DB{}, fmt.Errorf("failed to connect database")
	}
	return db, nil
}

func GetAllProperties(w http.ResponseWriter, r *http.Request) {
	db, err := GetDatabase()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	var properties []Property
	db.Table("Properties").Find(&properties)
	jsonBytes, err := json.Marshal(properties)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	w.Write(jsonBytes)
}

func AddProperty(w http.ResponseWriter, r *http.Request) {
	db, err := GetDatabase()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	var p Property
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	p.ID, err = uuid.NewUUID()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	result := db.Table("Properties").Create(&p)
	if result.Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(result.Error.Error()))
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Inserted Successfully"))
}
