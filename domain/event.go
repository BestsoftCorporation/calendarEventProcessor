package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event struct {
	ID            primitive.ObjectID `bson:"_id" json:"id"`
	CreatedDate   time.Time          `bson:"created_date" json:"created_date"`
	EndDate       time.Time          `bson:"end_date" json:"end_date"`
	StartDate     time.Time          `bson:"start_date" json:"start_date"`
	UpdateDate    time.Time          `bson:"update_date" json:"update_date"`
	EventId       string             `bson:"event_id" json:"event_id"`
	Location      string             `bson:"location" json:"location"`
	Mode          string             `bson:"mode" json:"mode"`
	Status        string             `bson:"status" json:"status"`
	Summary       string             `bson:"summary" json:"summary"`
	Attendees     []string           `bson:"attendees" json:"attendees"`
	Type          string             `bson:"type" json:"type"`
	UserEmail     string             `bson:"user_email" json:"user_email"`
	Processed     bool               `bson:"processed" json:"processed"`
	Commute       bool               `bson:"commute" json:"commute"`
	RemoteTime    time.Time          `bson:"remote_time" json:"remote_time"`
	RemoteTimeEnd time.Time          `bson:"remote_time_end" json:"remote_time_end"`
	BubbleId      string             `bson:"bubble_id" json:"bubble_id"`
}

type EventRepository interface {
	InsertOne(ctx context.Context, u *Event) (*Event, error)
	FindOne(ctx context.Context, id string) (*Event, error)
	FindOneByID(ctx context.Context, id string) (*Event, error)
	DeleteOneByID(ctx context.Context, id string) (*Event, error)
	Find(ctx context.Context, email string, date string) (*[]Event, error)
	FindCommute(ctx context.Context, date string) (*[]Event, error)
	UpdateOne(ctx context.Context, user *Event, id string) (*Event, error)
	DeleteAll(ctx context.Context, user_email string) (*[]Event, error)
}

type EventUsecase interface {
	InsertOne(ctx context.Context, u *Event) (*Event, error)
	FindOne(ctx context.Context, email string) (*Event, error)
	FindOneById(ctx context.Context, id string) (*Event, error)
	DeleteOneByID(ctx context.Context, id string) (*Event, error)
	Find(ctx context.Context, email string, date string) (*[]Event, error)
	FindCommute(ctx context.Context, date string) (*[]Event, error)
	UpdateOne(ctx context.Context, user *Event, id string) (*Event, error)
	DeleteAll(ctx context.Context, user_email string) (*[]Event, error)
}
