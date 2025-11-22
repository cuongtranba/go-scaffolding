package postgres

import (
	"time"

	"gorm.io/gorm"
)

// UserModel represents the database model for users
type UserModel struct {
	ID        string         `gorm:"type:uuid;primaryKey"`
	Email     string         `gorm:"type:varchar(254);uniqueIndex;not null"`
	Name      string         `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name for UserModel
func (UserModel) TableName() string {
	return "users"
}
