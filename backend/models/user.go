package models

import "gorm.io/gorm"

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Name     string
	Notes    []Note `gorm:"foreignKey:UserID"`
	gorm.Model
}