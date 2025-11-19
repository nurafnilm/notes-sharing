package models

import "gorm.io/gorm"

type Log struct {
	ID            uint   `gorm:"primaryKey"`
	Timestamp     string `json:"timestamp"`
	Method        string `json:"method"`
	Endpoint      string `json:"endpoint"`
	Headers       string `json:"headers"` // JSON string
	Payload       string `json:"payload,omitempty"`
	ResponseBody  string `json:"response_body"`
	StatusCode    int    `json:"status_code"`
	UserID        *uint  `json:"user_id,omitempty"`
	gorm.Model
}