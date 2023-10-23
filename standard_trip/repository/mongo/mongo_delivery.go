package mongo

import (
	"context"
	"github.com/bxcodec/go-clean-arch/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRepository struct {
	DB         mongo.Database
	Collection *mongo.Collection
	Some       string
}

const (
	timeFormat     = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
	collectionName = "recurringEvent"
)

func NewMongoRepository(DB mongo.Database) domain.StandardTripRepository {
	return &mongoRepository{DB, DB.Collection(collectionName), "sss"}
}

func (m *mongoRepository) InsertOne(ctx context.Context, recurringEvent *domain.StandardTrip) (*domain.StandardTrip, error) {
	var (
		err error
	)

	_, err = m.Collection.InsertOne(ctx, recurringEvent)

	if err != nil {
		return recurringEvent, err
	}

	return recurringEvent, nil
}

func (m *mongoRepository) FindAll(ctx context.Context) (*[]domain.StandardTrip, error) {
	var (
		recEvent []domain.StandardTrip
		err      error
	)

	cursor, err := m.Collection.Find(ctx, bson.D{})

	for cursor.Next(context.TODO()) {
		//Create a value into which the single document can beM decoded
		var elem domain.StandardTrip
		err := cursor.Decode(&elem)
		if err != nil {
			println(err)
		}
		recEvent = append(recEvent, elem)
	}

	if err != nil {
		println(err)
		return nil, err
	}

	return &recEvent, nil
}

func (m *mongoRepository) FindOne(ctx context.Context, userEmail string) (*domain.StandardTrip, error) {
	var (
		recEvent domain.StandardTrip
		err      error
	)

	println("Searching for " + userEmail)

	err = m.Collection.FindOne(ctx, bson.M{"user_email": userEmail}).Decode(&recEvent)
	if err != nil {
		return &recEvent, err
	}

	return &recEvent, nil
}

func (m *mongoRepository) DeleteOne(ctx context.Context, userEmail string) (*domain.StandardTrip, error) {
	var (
		recEvent domain.StandardTrip
		err      error
	)

	_, err = m.Collection.DeleteOne(ctx, bson.M{"user_email": userEmail})
	if err != nil {
		return &recEvent, err
	}

	return &recEvent, nil
}

func (m *mongoRepository) UpdateOne(ctx context.Context, standard *domain.StandardTrip, id string) (*domain.StandardTrip, error) {
	var (
		standEv domain.StandardTrip
		err     error
	)

	_, err = m.Collection.UpdateOne(ctx, bson.M{"userWorkspaceAddress": standard.UserWorkspaceAddress, "mode": standard.Mode}, bson.M{"user_email": id})
	if err != nil {
		return &standEv, err
	}

	return &standEv, nil
}
