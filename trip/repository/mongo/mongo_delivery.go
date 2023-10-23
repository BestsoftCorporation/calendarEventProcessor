package mongo

import (
	"context"
	"github.com/bxcodec/go-clean-arch/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRepository struct {
	DB         mongo.Database
	Collection *mongo.Collection
}

const (
	timeFormat     = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
	collectionName = "Trip"
)

func NewMongoRepository(DB mongo.Database) domain.TripRepository {
	return &mongoRepository{DB, DB.Collection(collectionName)}
}

func (m *mongoRepository) InsertOne(ctx context.Context, Trip *domain.Trip) (*domain.Trip, error) {
	var (
		err error
	)

	_, err = m.Collection.InsertOne(ctx, Trip)

	if err != nil {
		return Trip, err
	}

	return Trip, nil
}

func (m *mongoRepository) FindOne(ctx context.Context, id string, returning bool) (*domain.Trip, error) {
	var (
		trip domain.Trip
		err  error
	)

	err = m.Collection.FindOne(ctx, bson.M{"linkedID": id, "returning": returning}).Decode(&trip)
	if err != nil {
		return &trip, err
	}

	return &trip, nil
}

func (m *mongoRepository) DeleteAll(ctx context.Context, userEmail string) (*domain.Trip, error) {
	var (
		trip domain.Trip
		err  error
	)

	_, err = m.Collection.DeleteMany(ctx, bson.M{"userEmail": userEmail})
	if err != nil {
		return &trip, err
	}

	return &trip, nil
}

func (m *mongoRepository) FindAll(ctx context.Context, userEmail string) (*[]domain.Trip, error) {
	var (
		trips []domain.Trip
		err   error
	)

	cursor, err := m.Collection.Find(ctx, bson.M{"userEmail": userEmail})

	for cursor.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem domain.Trip
		err := cursor.Decode(&elem)
		if err != nil {
			println(err)
		}
		trips = append(trips, elem)
	}

	if err != nil {
		println(err)
		return nil, err
	}

	return &trips, nil
}

func (m *mongoRepository) UpdateOne(ctx context.Context, trip *domain.Trip, id primitive.ObjectID) (*domain.Trip, error) {
	var (
		tripOne domain.Trip
		err     error
	)

	//{"destinationGeo", trip.DestinationGeo}, {"startLocationGeo", trip.StartLocationGeo},
	_, err = m.Collection.UpdateOne(ctx, bson.D{{"_id", id}}, bson.D{{"$set",
		bson.D{
			{"startLocation", trip.StartLocation},
			{"remote_hours", trip.RemoteHours},
			{"destination", trip.Destination},
			{"mode", trip.Mode},
			{"allowance_amount", trip.AllowanceAmount}}}})

	if err != nil {
		println(err.Error())
		return &tripOne, err
	}

	println("UPDATING...")

	return &tripOne, nil
}

func (m *mongoRepository) DeleteByReturningAndID(ctx context.Context, ret bool, id string, field string,
	disableID string, delete bool) (*domain.Trip, error) {
	var (
		trip domain.Trip
		err  error
	)

	m.Collection.FindOne(ctx, bson.D{{"returning", ret}, {"linkedID", id}}).Decode(&trip)
	if delete {
		_, err = m.Collection.DeleteOne(ctx, bson.D{{"returning", ret}, {"linkedID", id}})
	} else {
		_, err = m.Collection.UpdateOne(ctx, bson.D{{"returning", ret}, {"linkedID", id}},
			bson.D{{"$set", bson.D{{"disabled", true}, {"bubble_id", ""},
				{"disabledId", disableID}}}})
	}

	if err != nil {
		return &trip, err
	}
	println("Trying to delete with id " + id + " and field" + field)

	return &trip, nil
}

func (m *mongoRepository) DeleteAllWhereLinkedId(ctx context.Context, linkedId string) (*domain.Trip, error) {
	var (
		trip domain.Trip
		err  error
	)

	_, err = m.Collection.UpdateMany(ctx, bson.D{{"disabledId", linkedId}}, bson.D{{"$set", bson.D{{"disabled", false}, {"disabledId", ""}}}})
	if err != nil {
		return &trip, err
	}
	_, err = m.Collection.DeleteMany(ctx, bson.D{{"linkedID", linkedId}})
	if err != nil {
		return &trip, err
	}

	println("Trying to delete " + linkedId)

	return &trip, nil
}
