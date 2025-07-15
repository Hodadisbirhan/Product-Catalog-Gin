package model

type Role struct {
	BaseModel
	Name        string       `gorm:"unique;not null"`
	Permissions []Permission `gorm:"many2many:role_permissions"`
}
