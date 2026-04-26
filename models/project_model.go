package models

import (
	"time"

	"gorm.io/datatypes"
)

// ID -> Unique Id
// CreatedAt
// UpdatedAt
// Name
// Description
// Image
// Skills

type Project struct {
	ID          uint `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string `binding:"required"`
	Description string `binding:"required"`
	Image       string
	Skills      datatypes.JSONSlice[string] `gorm:"type:json" swaggertype:"array,string"`
	Comments    []Comment                   `gorm:"foreignKey:ProjectID"`
	Likes       []User                      `gorm:"many2many:project_likes"`
}

type ProjectUpdateInput struct {
	Name        *string   `json:"name"`
	Description *string   `json:"description"`
	Skills      *[]string `json:"skills"`
}
