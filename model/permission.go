package model

import "github.com/google/uuid"

type Permission struct {
	BaseModel
	Name          string       `gorm:"not null"`
	ContentTypeID uuid.UUID    `gorm:"type:uuid"`
	ContentType   ContentType  `gorm:"foreignKey:ContentTypeID"`
}
