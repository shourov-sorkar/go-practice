package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string             `json:"name" gorm:"not null"`
	Username  string             `json:"username" gorm:"unique;not null"`
	Email     string             `json:"email" gorm:"unique;not null"`
	Password  string             `json:"password" gorm:"not null"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}
