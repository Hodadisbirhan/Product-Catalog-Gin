package model

import "github.com/google/uuid"

type User struct {
	BaseModel
	Username string `gorm:"unique;not null"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`

	RoleID uuid.UUID `gorm:"type:uuid"`
	Role   Role
}
