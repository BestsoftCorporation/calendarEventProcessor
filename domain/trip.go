package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Trip struct {
	ID               primitive.ObjectID `bson:"_id" json:"id"`
	Date             time.Time          `bson:"date" json:"date"`
	CreatedDate      time.Time          `bson:"created_date" json:"created_date"`
	UpdateDate       time.Time          `bson:"update_date" json:"update_date"`
	StartTime        time.Time          `bson:"startTime" json:"startTime"`
	AllowanceAmount  int                `bson:"allowance_amount" json:"allowance_amount"`
	StartLocation    string             `bson:"startLocation" json:"startLocation"`
	StartLocationGeo Geo                `bson:"startLocationGeo" json:"startLocationGeo"`
	Destination      string             `bson:"destination" json:"destination"`
	DestinationGeo   Geo                `bson:"destinationGeo" json:"destinationGeo"`
	Summary          string             `bson:"summary" json:"summary"`
	Disabled         bool               `bson:"disabled" json:"disabled"`
	DisabledID       string             `bson:"disabledId" json:"disabledId"`
	LinkedID         string             `bson:"linkedID" json:"linkedID"`
	Mode             string             `bson:"mode" json:"mode"`
	Commute          bool               `bson:"commute" json:"commute"`
	RemoteDay        bool               `bson:"remoteDay" json:"remoteDay"`
	UserEmail        string             `bson:"userEmail" json:"userEmail"`
	Validate         bool               `bson:"validate" json:"validate"`
	Returning        bool               `bson:"returning" json:"returning"`
	BubbleId         string             `bson:"bubble_id" json:"bubble_id"`
	RemoteHours      int                `bson:"remote_hours" json:"remote_hours"`
}

type Geo struct {
	Lat float64
	Log float64
}

type TripRepository interface {
	InsertOne(ctx context.Context, u *Trip) (*Trip, error)
	FindOne(ctx context.Context, id string, returning bool) (*Trip, error)
	FindAll(ctx context.Context, userEmail string) (*[]Trip, error)
	UpdateOne(ctx context.Context, user *Trip, id primitive.ObjectID) (*Trip, error)
	DeleteAll(ctx context.Context, userEmail string) (*Trip, error)
	DeleteByReturningAndID(ctx context.Context, ret bool, id string, field string, disableID string, delete bool) (*Trip,
		error)
	DeleteAllWhereLinkedId(ctx context.Context, linkedId string) (*Trip, error)
}

type TripUsecase interface {
	InsertOne(ctx context.Context, u *Trip) (*Trip, error)
	FindOne(ctx context.Context, id string, returning bool) (*Trip, error)
	FindAll(ctx context.Context, userEmail string) (*[]Trip, error)
	UpdateOne(ctx context.Context, user *Trip, id primitive.ObjectID) (*Trip, error)
	DeleteAll(ctx context.Context, userEmail string) (*Trip, error)
	DeleteByReturningAndID(ctx context.Context, ret bool, id string, field string, disableID string, delete bool) (*Trip,
		error)
	DeleteAllWhereLinkedId(ctx context.Context, linkedId string) (*Trip, error)
}
