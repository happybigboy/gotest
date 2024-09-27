// models/models.go
package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ChatID   int64  `gorm:"uniqueIndex"`
	Username string `gorm:"not null"` // Ensure username is not null
	Password string `gorm:"not null"` // Ensure password is not null
	Token 	 string `gorm:"not null"`
}

// Function to create or update a user
func CreateUser(db *gorm.DB, chatID int64, username, password string,token string) error {
	user := User{
		ChatID:   chatID,
		Username: username,
		Password: password,
		Token: token,
	}

	// Check if the user already exists
	var existingUser User
	if err := db.Where("chat_id = ?", chatID).First(&existingUser).Error; err == nil {
		// User exists, update the existing record
		existingUser.Username = username
		existingUser.Password = password
		return db.Save(&existingUser).Error
	}

	// User does not exist, create a new record
	return db.Create(&user).Error
}
func ReadUser(db *gorm.DB, chatID int64) (*User, error) {
	var user User
	if err := db.Where("chat_id = ?", chatID).First(&user).Error; err != nil {
		return nil, err // Return nil and the error if the user does not exist
	}
	return &user, nil // Return the user if found
}
func ModifyUser(db *gorm.DB, chatID int64, username, password string) error {
	var user User
	if err := db.Where("chat_id = ?", chatID).First(&user).Error; err != nil {
		return err // Return the error if the user does not exist
	}

	// Update user fields
	user.Username = username
	user.Password = password
	return db.Save(&user).Error // Save the modified user
}