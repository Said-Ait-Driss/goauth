package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id"`
	First_name    *string            `json:"first_name" validate:"required, min=3, max=100"`
	Last_name     *string            `json:"last_name" validate:"required, min=3, max=100"`
	Phone         *string            `json:"phone" validate:"required, min=10, max=10"`
	Email         *string            `json:"email" validate:"required, email"`
	Password      *string            `json:"password" validate:"min=6"`
	Created_at    time.Time          `json:"created_at"`
	Updated_at    time.Time          `json:"updated_at"`
	Token         *string            `json:"token"`
	Refresh_token *string            `json:"refresh_token"`
	User_type     *string            `json:"user_type" validate:"required, eq=admin|eq=user"`
	User_id       string             `json:"user_id"`
}
