package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Services struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	ServiceName  string             `bson:"ServiceName" json:"ServiceName" validate:"required"`
	ServiceEmail string             `bson:"ServiceEmail" json:"ServiceEmail" validate:"required"`
	ServicePrice float32            `bson:"ServicePrice" json:"ServicePrice" validate:"required"`
}
