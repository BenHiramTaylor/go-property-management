package users

import (
	"fmt"
	"net/http"

	"github.com/BenHiramTaylor/go-property-management/database"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model `json:"-"`
	ID         uuid.UUID `json:"id" validate:"required"`
	Username   string    `json:"username" validate:"required"`
	FirstName  string    `json:"first_name" validate:"required"`
	MiddleName string    `json:"middle_name,omitempty"`
	LastName   string    `json:"last_name" validate:"required"`
	Email      string    `json:"email" validate:"required,email"`
	Password   string    `json:"password" validate:"required,min=8"`
}

type UserResonse struct {
	ID         uuid.UUID `json:"id"`
	Username   string    `json:"username"`
	FirstName  string    `json:"first_name"`
	MiddleName string    `json:"middle_name,omitempty"`
	LastName   string    `json:"last_name"`
	Email      string    `json:"email"`
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (u *User) validate() error {
	v := validator.New()
	return v.Struct(u)
}

func (u *User) hashPassword() error {
	bytesHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {
		return err
	}
	hashedPassword := string(bytesHash)
	isCorrect := checkPasswordHash(u.Password, hashedPassword)
	if !isCorrect {
		return fmt.Errorf("hashed password does not match raw password")
	}
	u.Password = hashedPassword
	return nil
}

func AddUser(c *fiber.Ctx) error {
	var newUser User
	err := c.BodyParser(&newUser)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	newUser.ID, err = uuid.NewUUID()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	err = newUser.validate()
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	err = newUser.hashPassword()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	result := database.DBConn.Table("Users").Create(&newUser)
	if result.Error != nil {
		return result.Error
	}
	userResponse := UserResonse{
		ID:         newUser.ID,
		Username:   newUser.Username,
		FirstName:  newUser.FirstName,
		MiddleName: newUser.MiddleName,
		LastName:   newUser.LastName,
		Email:      newUser.Email}
	return c.JSON(userResponse)
}
