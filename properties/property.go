package properties

import (
	"fmt"
	"log"
	"net/http"

	"github.com/BenHiramTaylor/go-property-management/database"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Property struct {
	gorm.Model          `json:"-"`
	ID                  uuid.UUID `json:"id"`
	PropertyType        string    `json:"property_type"`
	Address             string    `json:"address"`
	NumberOfBedrooms    uint      `json:"number_of_bedrooms"`
	PurchasePriceDollar uint      `json:"purchase_price_dollar"`
}

func GetIndividualPropertyByID(id string) (*Property, error) {
	var p Property
	result := database.DBConn.Table("Properties").Find(&p, "id = ?", id)
	if result.Error != nil {
		return &p, result.Error
	}
	if result.RowsAffected == 0 {
		return &p, fmt.Errorf("property not found with ID: %v", id)
	}
	return &p, nil
}

func GetAllProperties(c *fiber.Ctx) error {
	var properties []Property
	result := database.DBConn.Table("Properties").Find(&properties)
	if result.Error != nil {
		return result.Error
	}
	return c.JSON(properties)
}

func AddProperty(c *fiber.Ctx) error {
	var p Property
	err := c.BodyParser(&p)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}
	log.Println(fmt.Sprintf("%#v", p))
	p.ID, err = uuid.NewUUID()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}
	result := database.DBConn.Table("Properties").Create(&p)
	if result.Error != nil {
		return result.Error
	}
	return c.JSON(p)
}

func GetIndividualProperty(c *fiber.Ctx) error {
	id := c.Params("id")
	p, err := GetIndividualPropertyByID(id)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(err.Error())
	}
	return c.JSON(p)
}

func DeleteProperty(c *fiber.Ctx) error {
	id := c.Params("id")
	p, err := GetIndividualPropertyByID(id)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(err.Error())
	}
	database.DBConn.Table("Properties").Delete(&p)
	return c.JSON("Property Successfully Deleted")
}

func UpdateProperty(c *fiber.Ctx) error {
	id := c.Params("id")
	var newP Property
	oldP, err := GetIndividualPropertyByID(id)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(err.Error())
	}
	err = c.BodyParser(&newP)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}
	// SET NEW ID TO THE ID FROM THE URL
	newP.ID = oldP.ID
	result := database.DBConn.Table("Properties").Model(&oldP).Updates(newP)
	if result.Error != nil {
		return result.Error
	}
	return c.JSON(&newP)
}
