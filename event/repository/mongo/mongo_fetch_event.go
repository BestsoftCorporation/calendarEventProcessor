package mongo

import (
	"context"

	"github.com/bxcodec/go-clean-arch/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type fetchEventmongoRepository struct {
	DB         mongo.Database
	Collection *mongo.Collection
}

const (
	collectionName2 = "FetchEvent"
)

func NewFetchEventMongoRepository(DB mongo.Database) domain.EventFetchRepository {
	return &fetchEventmongoRepository{DB, DB.Collection(collectionName2)}
}

func (m *fetchEventmongoRepository) InsertOne(ctx context.Context, Event *domain.EventFetch) (*domain.EventFetch, error) {
	var (
		err error
	)

	_, err = m.Collection.InsertOne(ctx, Event)

	if err != nil {
		return Event, err
	}

	return Event, nil
}

func (m *fetchEventmongoRepository) FindOne(ctx context.Context, id string) (*domain.EventFetch, error) {
	var (
		fetchEvent domain.EventFetch
		err        error
	)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &fetchEvent, err
	}

	err = m.Collection.FindOne(ctx, bson.M{"_id": idHex}).Decode(&fetchEvent)
	if err != nil {
		return &fetchEvent, err
	}

	return &fetchEvent, nil
}

func (m *fetchEventmongoRepository) FindOneByType(ctx context.Context, typeOfEvent string) (*domain.EventFetch, error) {
	var (
		fetchEvent domain.EventFetch
		err        error
	)

	err = m.Collection.FindOne(ctx, bson.M{"type": typeOfEvent}).Decode(&fetchEvent)
	if err != nil {
		return &fetchEvent, err
	}

	return &fetchEvent, nil
}

func (m *fetchEventmongoRepository) UpdateOne(ctx context.Context, ev *domain.EventFetch, eventID string) (*domain.EventFetch, error) {

	var (
		fetchEvent domain.EventFetch
		err        error
	)

	idHex, err := primitive.ObjectIDFromHex(eventID)
	if err != nil {
		return &fetchEvent, err
	}

	_, err = m.Collection.UpdateOne(ctx, bson.D{{"_id", idHex}}, bson.D{{"$set", bson.D{{"last_fetch", ev.LastFetch}}}})
	if err != nil {
		println(err.Error())
		return &fetchEvent, err
	}
	println("UPDATING...")

	return &fetchEvent, nil
}
