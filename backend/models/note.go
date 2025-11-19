package models

import "gorm.io/gorm"

type Note struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	ImageURL string `json:"image_url,omitempty" gorm:"column:image_url"`
	UserID   uint   `json:"-"` // hidden dari response
	gorm.Model
}