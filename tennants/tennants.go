package tennants

import (
	"fmt"
	"net/http"
	"time"

	"github.com/BenHiramTaylor/go-property-management/database"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tennant struct {
	gorm.Model
	ID          uuid.UUID `json:"id"`
	FirstName   string    `json:"first_name"`
	MiddleName  string    `json:"middle_name"`
	LastName    string    `json:"last_name"`
	DateOfBirth time.Time `json:"date_of_birth"`
	PropertyID  uuid.UUID `json:"property_id"`
}

func GetAllTennants(c *fiber.Ctx) error {
	var tennants []Tennant
	result := database.DBConn.Table("Properties").Find(&tennants)
	if result.Error != nil {
		return result.Error
	}
	return c.JSON(tennants)
}

func AddTennant(c *fiber.Ctx) error {
	var t Tennant
	err := c.BodyParser(&t)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}
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
	var t Tennant
	result := database.DBConn.Table("Tennants").Find(&t, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return c.Status(http.StatusNotFound).JSON(fmt.Sprintf("Tennant not found with id: %v", id))
	}
	return c.JSON(t)
}

func DeleteTennant(c *fiber.Ctx) error {
	id := c.Params("id")
	var t Tennant
	result := database.DBConn.Table("Tennants").Find(&t, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return c.Status(http.StatusNotFound).JSON(fmt.Sprintf("Tennant not found with id: %v", id))
	}
	database.DBConn.Table("Tennants").Delete(&t)
	return c.JSON("Tennant Successfully Deleted")
}
