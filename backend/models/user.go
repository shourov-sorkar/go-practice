package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" binding:"required" bson:"name"`
	Username  string             `json:"username" binding:"required" bson:"username"`
	Email     string             `json:"email" binding:"required,email" bson:"email"`
	Password  string             `json:"password" binding:"required,min=6" bson:"password"`
	CreatedAt time.Time          `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updated_at"`
}
