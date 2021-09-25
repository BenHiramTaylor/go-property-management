package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/BenHiramTaylor/go-property-management/properties"
	"github.com/BenHiramTaylor/go-property-management/tennants"
	"github.com/gorilla/mux"
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

func InitialiseDB() error {
	db, err := GetDatabase()
	if err != nil {
		return err
	}
	err = db.Table("Properties").AutoMigrate(&properties.Property{})
	if err != nil {
		return err
	}
	err = db.Table("Tennants").AutoMigrate(&tennants.Tennant{})
	if err != nil {
		return err
	}
	return nil
}

func initialiseRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/properties", properties.GetAllProperties).Methods(http.MethodGet)
	r.HandleFunc("/properties", properties.AddProperty).Methods(http.MethodPost)
	return r
}

func main() {
	log.Println("Starting Server...")
	err := InitialiseDB()
	if err != nil {
		panic(err)
	}
	log.Println("Initialised DB")
	r := initialiseRoutes()
	srv := &http.Server{Addr: ":80", Handler: r, WriteTimeout: 15 * time.Second, ReadTimeout: 15 * time.Second}
	log.Fatalln(srv.ListenAndServe())
}
