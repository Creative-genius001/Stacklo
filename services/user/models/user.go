package models

import "time"

type User struct {
	Id           string `gorm:"type:uuid;primaryKey;unique"`
	Email        string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
	FirstName    string
	LastName     string
	Phone        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
