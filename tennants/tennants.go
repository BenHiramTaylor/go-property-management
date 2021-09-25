package tennants

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tennant struct {
	gorm.Model
	ID          uuid.UUID
	FirstName   string
	MiddleName  string
	LastName    string
	DateOfBirth time.Time
	PropertyID  uuid.UUID
}
