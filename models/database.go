// models/models.go
package models

import (
	"gorm.io/gorm"
)

// StateModel represents the 'states' table
type StateModel struct {
	UserID int    `gorm:"primaryKey"`
	ChatID int    `gorm:"primaryKey"`
	State  string `gorm:"type:varchar(255)"`
}

// DataModel represents the 'data' table
type DataModel struct {
	UserID int         `gorm:"primaryKey"`
	ChatID int         `gorm:"primaryKey"`
	Data   interface{} `gorm:"type:json"`
}

// MessageModel represents the 'messages' table
type MessageModel struct {
	ID        int `gorm:"primaryKey;autoIncrement"`
	ChatID    int `gorm:"index"`
	MessageID int `gorm:"index"`
}

// MigrateModels migrates the tables in the provided DB
func MigrateModels(db *gorm.DB) error {
	// AutoMigrate the models
	return db.AutoMigrate(&StateModel{}, &DataModel{}, &MessageModel{})
}

// Custom exceptions
type CustomException struct {
	Message string
}

func (e *CustomException) Error() string {
	return e.Message
}

type AuthenticationError struct {
	CustomException
}

type UserNotFoundError struct {
	CustomException
}

type APIError struct {
	CustomException
}

type NetworkError struct {
	CustomException
}

// Custom exception helpers
func NewAuthenticationError() error {
	return &AuthenticationError{CustomException{Message: "Authentication failed"}}
}

func NewUserNotFoundError() error {
	return &UserNotFoundError{CustomException{Message: "User not found"}}
}

func NewAPIError() error {
	return &APIError{CustomException{Message: "API error occurred"}}
}

func NewNetworkError() error {
	return &NetworkError{CustomException{Message: "Network error occurred"}}
}