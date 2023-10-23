package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Maps struct {
	ID     primitive.ObjectID `bson:"_id" json:"id"`
	APIKey string             `bson:"APIKey" json:"APIKey" validate:"required"`
}
