package models

import "gorm.io/gorm"

type User struct {
	gorm.Model        // Embedding Gorm's base model which includes ID, CreatedAt, UpdatedAt, DeletedAt
	Email      string `gorm:"type:varchar(100);uniqueIndex"` // User's email with a unique index
	Password   string // User's hashed password
}
