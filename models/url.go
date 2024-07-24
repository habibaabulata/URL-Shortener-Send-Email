package models

import "gorm.io/gorm"

type URL struct {
    gorm.Model
    ShortCode   string `gorm:"uniqueIndex"`
    OriginalURL string
    UserID      uint
}
