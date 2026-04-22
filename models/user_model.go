package models

import "time"

type User struct {
	Id        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Email     string `gorm:"unique" binding:"required,email"`
	Password  string `binding:"required,min=6"`
}
