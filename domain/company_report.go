package domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type CompanyReport struct {
	ID                primitive.ObjectID `bson:"_id" json:"id"`
	From              time.Time          `bson:"from" json:"from"`
	CreatedAt         time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt         time.Time          `bson:"updated_at" json:"updated_at"`
	File_url          string             `json:"file_url" `
	Year              string             `json:"year"`
	Month             string             `json:"month" `
	Spent             string             `json:"spent"`
	Remote_days_spent string             `json:"remote_days_spent" `
}

type CompanyReportRepository interface {
	InsertOne(ctx context.Context, u *CompanyReport) (*CompanyReport, error)
	FindOne(ctx context.Context, id string) (*CompanyReport, error)
	UpdateOne(ctx context.Context, user *CompanyReport, id string) (*CompanyReport, error)
}

type CompanyReportUsecase interface {
	InsertOne(ctx context.Context, u *CompanyReport) (*CompanyReport, error)
	FindOne(ctx context.Context, id string) (*CompanyReport, error)
	UpdateOne(ctx context.Context, user *CompanyReport, id string) (*CompanyReport, error)
}
