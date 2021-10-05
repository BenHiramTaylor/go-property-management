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

// User The primary user struct, contains all the user information including gorm.Model and various struct tags
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

// UserResponse Is a struct simply for returning different user fields after a successful reqeust,
// to prevent passwords from being passed over http
type UserResponse struct {
	ID         uuid.UUID `json:"id"`
	Username   string    `json:"username"`
	FirstName  string    `json:"first_name"`
	MiddleName string    `json:"middle_name,omitempty"`
	LastName   string    `json:"last_name"`
	Email      string    `json:"email"`
}

func (u *User) createUserResponse() *UserResponse {
	return &UserResponse{ID: u.ID, Username: u.Username, FirstName: u.FirstName, MiddleName: u.MiddleName, LastName: u.LastName, Email: u.Email}
}

// checkPasswordHash Takes a password, and a hash string, and compares the two, if they do not match then return false
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// validate Checks the struct tags on the User object to validate the properties, returns error if validation fails
func (u *User) validate() error {
	v := validator.New()
	return v.Struct(u)
}

// hashPassword A method that hashes the plain text password of a new user struct, and then replaces the
// plain text password with the hashed string
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

// getUserByUsername Gets a user from the DB table via the username key, or returns an error
func getUserByUsername(username string) (*User, error) {
	var u User
	result := database.DBConn.Table("Users").Find(&u, "username = ?", username)
	if result.Error != nil {
		return &u, result.Error
	}
	if result.RowsAffected == 0 {
		return &u, fmt.Errorf("username not found: %v", username)
	}
	return &u, nil
}

// checkIfUserExists Returns true if a user exists with the given username, and false if not
func checkIfUserExists(username string) bool {
	var u User
	result := database.DBConn.Table("Users").Find(&u, "username = ?", username)
	if result.RowsAffected == 0 {
		return false
	}
	return true
}

// CheckAuth Takes a username and password, and compares the password to the hash stored in the Users table to
// authenticate, returns true if authenticated and false if not
func CheckAuth(username, password string) bool {
	user, err := getUserByUsername(username)
	if err != nil {
		return false
	}
	passwordMatch := checkPasswordHash(password, user.Password)
	if !passwordMatch {
		return false
	}
	return true
}

// AddUser Takes a fiber context, and creates a new User in the DB via a post request, or returns an error to the client
func AddUser(c *fiber.Ctx) error {
	var newUser User
	err := c.BodyParser(&newUser)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	userExists := checkIfUserExists(newUser.Username)
	if userExists {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "username provided already exists"})
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
	userResponse := newUser.createUserResponse()
	return c.JSON(*userResponse)
}
