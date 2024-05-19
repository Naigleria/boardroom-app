package models

import (
	"gorm.io/gorm"
)

//modelo para los usuarios
type User struct {
	gorm.Model
	Username          string `gorm:"unique"`
	Email             string `gorm:"unique"`
	Password          string
	VerificationToken string
	IsVerified        bool
}

type Users []User