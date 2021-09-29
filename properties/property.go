package properties

import (
	"fmt"
	"log"
	"net/http"

	"github.com/BenHiramTaylor/go-property-management/database"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Property struct {
	gorm.Model       `json:"-"`
	ID               uuid.UUID `json:"id" validate:"required"`
	PropertyType     string    `json:"property_type" validate:"required"`
	Address          string    `json:"address" validate:"required"`
	NumberOfBedrooms uint      `json:"number_of_bedrooms" validate:"required"`
	PurchasePriceGBP uint      `json:"purchase_price_gbp" validate:"required"`
}

func (p *Property) validate() error {
	v := validator.New()
	return v.Struct(p)
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
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	log.Println(fmt.Sprintf("%#v", p))
	p.ID, err = uuid.NewUUID()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	err = p.validate()
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
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
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(p)
}

func DeleteProperty(c *fiber.Ctx) error {
	id := c.Params("id")
	p, err := GetIndividualPropertyByID(id)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	database.DBConn.Table("Properties").Delete(&p)
	return c.JSON("Property Successfully Deleted")
}

func UpdateProperty(c *fiber.Ctx) error {
	id := c.Params("id")
	var newP Property
	oldP, err := GetIndividualPropertyByID(id)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	err = c.BodyParser(&newP)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	// SET NEW ID TO THE ID FROM THE URL
	newP.ID = oldP.ID
	result := database.DBConn.Table("Properties").Model(&oldP).Updates(newP)
	if result.Error != nil {
		return result.Error
	}
	return c.JSON(&newP)
}
