package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StandardTrip struct {
	ID                   primitive.ObjectID   `bson:"_id" json:"id"`
	From                 time.Time            `bson:"from" json:"from"`
	CreatedAt            time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt            time.Time            `bson:"updated_at" json:"updated_at"`
	UserEmail            string               `bson:"user_email" json:"user_email"`
	UserHomeAddress      string               `bson:"userHomeAddress" json:"userHomeAddress" `
	UserWorkspaceAddress string               `bson:"userWorkspaceAddress" json:"userWorkspaceAddress" `
	Mode                 string               `bson:"mode" json:"mode" `
	WeeklyDistance       int                  `bson:"weeklyDistance" json:"weeklyDistance" `
	Days                 []Day                `bson:"days" json:"days"`
	UserAnnualBudget     primitive.Decimal128 `bson:"userAnnualBudget" json:"userAnnualBudget" `
}

type Day struct {
	Distance int       `bson:"distance" json:"distance"`
	Trips    []DayTrip `bson:"trips" json:"trips"`
	Cost     int64     `bson:"cost" json:"cost"`
	Mode     string    `bson:"mode" json:"mode"`
	Day      string    `bson:"day" json:"day"`
}

type TripDetails struct {
	WorkType string `json:"workType"`
	Mode     string `json:"Mode"`
}

type DayTrip struct {
	Distance      int       `bson:"distance" json:"distance"`
	StartLocation string    `bson:"startLocation" json:"startLocation"`
	EndLocation   string    `bson:"endLocation" json:"endLocation"`
	StartTime     time.Time `bson:"startTime" json:"startTime"`
	EndTime       time.Time `bson:"endTime" json:"endTime"`
	Service       string    `bson:"service" json:"service"`
}

type StandardTripRepository interface {
	InsertOne(ctx context.Context, u *StandardTrip) (*StandardTrip, error)
	FindOne(ctx context.Context, id string) (*StandardTrip, error)
	FindAll(ctx context.Context) (*[]StandardTrip, error)
	UpdateOne(ctx context.Context, user *StandardTrip, id string) (*StandardTrip, error)
	DeleteOne(ctx context.Context, email string) (*StandardTrip, error)
}

type StandardTripUsecase interface {
	InsertOne(ctx context.Context, u *StandardTrip) (*StandardTrip, error)
	FindOne(ctx context.Context, id string) (*StandardTrip, error)
	FindAll(ctx context.Context) (*[]StandardTrip, error)
	UpdateOne(ctx context.Context, user *StandardTrip, id string) (*StandardTrip, error)
	DeleteOne(ctx context.Context, email string) (*StandardTrip, error)
}
