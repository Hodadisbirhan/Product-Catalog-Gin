package model

type ContentType struct {
	BaseModel
	Name string `gorm:"unique;not null"`
}
