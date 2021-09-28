package tennants

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/BenHiramTaylor/go-property-management/database"
	"github.com/BenHiramTaylor/go-property-management/properties"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tennant struct {
	gorm.Model  `json:"-"`
	ID          uuid.UUID `json:"id"`
	FirstName   string    `json:"first_name"`
	MiddleName  string    `json:"middle_name"`
	LastName    string    `json:"last_name"`
	DateOfBirth time.Time `json:"date_of_birth"`
	PropertyID  uuid.UUID `json:"property_id"`
}

func GetIndividualTennantByID(id string) (*Tennant, error) {
	var t Tennant
	result := database.DBConn.Table("Tennants").Find(&t, "id = ?", id)
	if result.Error != nil {
		return &t, result.Error
	}
	if result.RowsAffected == 0 {
		return &t, fmt.Errorf("tennant not found with ID: %v", id)
	}
	return &t, nil
}

func GetAllTennants(c *fiber.Ctx) error {
	var tennants []Tennant
	result := database.DBConn.Table("Tennants").Find(&tennants)
	if result.Error != nil {
		return result.Error
	}
	return c.JSON(tennants)
}

func AddTennant(c *fiber.Ctx) error {
	var t Tennant
	err := json.Unmarshal(c.Body(), &t)
	if err != nil {
		log.Println(err.Error())
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}
	log.Println(fmt.Sprintf("%#v", t))
	t.ID, err = uuid.NewUUID()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}
	result := database.DBConn.Table("Tennants").Create(&t)
	if result.Error != nil {
		return result.Error
	}
	return c.JSON(&t)
}

func GetIndividualTennant(c *fiber.Ctx) error {
	id := c.Params("id")
	t, err := GetIndividualTennantByID(id)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(err.Error())
	}
	return c.JSON(t)
}

func DeleteTennant(c *fiber.Ctx) error {
	id := c.Params("id")
	t, err := GetIndividualTennantByID(id)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(err.Error())
	}
	database.DBConn.Table("Tennants").Delete(&t)
	return c.JSON("Tennant Successfully Deleted")
}

func UpdateTennant(c *fiber.Ctx) error {
	id := c.Params("id")
	var newT Tennant
	oldT, err := GetIndividualTennantByID(id)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(err.Error())
	}
	err = c.BodyParser(&newT)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}
	// SET NEW ID TO THE ID FROM THE URL
	newT.ID = oldT.ID
	result := database.DBConn.Table("Tennants").Model(&oldT).Updates(newT)
	if result.Error != nil {
		return result.Error
	}
	return c.JSON(&newT)
}

func AssignTennantToProperty(c *fiber.Ctx) error {
	tennantID := c.Params("tennantID")
	propertyID := c.Params("propertyID")
	p, err := properties.GetIndividualPropertyByID(propertyID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(err.Error())
	}
	t, err := GetIndividualTennantByID(tennantID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(err.Error())
	}
	result := database.DBConn.Table("Tennants").Model(&t).Update("PropertyID", p.ID)
	if result.Error != nil {
		return result.Error
	}
	return c.JSON(&t)
}
