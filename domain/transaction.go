package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	ID                   primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt            time.Time          `bson:"created_at" json:"created_at"`
	UserID               primitive.ObjectID `bson:"_id" json:"user_id"`
	UserHomeAddress      string             `bson:"userHomeAddress" json:"userHomeAddress" validate:"required"`
	UserWorkspaceAddress string             `bson:"userWorkspaceAddress" json:"userWorkspaceAddress" validate:"required"`

	Cost primitive.Decimal128 `bson:"Cost" json:"Cost" validate:"required"`
}
