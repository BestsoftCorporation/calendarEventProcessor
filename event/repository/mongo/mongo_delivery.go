package mongo

import (
	"context"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

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
	collectionName = "Event"
)

func NewMongoRepository(DB mongo.Database) domain.EventRepository {
	return &mongoRepository{DB, DB.Collection(collectionName)}
}

func (m *mongoRepository) InsertOne(ctx context.Context, Event *domain.Event) (*domain.Event, error) {
	var (
		err error
	)

	_, err = m.Collection.InsertOne(ctx, Event)

	if err != nil {
		return Event, err
	}

	return Event, nil
}

func (m *mongoRepository) FindOne(ctx context.Context, userEmail string) (*domain.Event, error) {
	var (
		recEvent domain.Event
		err      error
	)

	println("Searching for " + userEmail)

	err = m.Collection.FindOne(ctx, bson.M{"user_email": userEmail}).Decode(&recEvent)
	if err != nil {
		return &recEvent, err
	}

	return &recEvent, nil
}

func (m *mongoRepository) Find(ctx context.Context, userEmail string, date string) (*[]domain.Event, error) {
	var (
		recEvent []domain.Event
		err      error
	)

	t, _ := time.Parse(time.RFC3339, strings.Split(date, "T")[0]+"T00:00:00.000+00:00")

	opts := options.Find().SetSort(bson.D{{"start_date", 1}})

	cursor, err := m.Collection.Find(ctx, bson.D{{"user_email", userEmail}, {"start_date", bson.M{"$gte": primitive.NewDateTimeFromTime(t)}}, {"start_date", bson.M{"$lt": primitive.NewDateTimeFromTime(t.AddDate(0, 0, 1))}}}, opts)

	for cursor.Next(context.TODO()) {
		//Create a value into which the single document can beM decoded
		var elem domain.Event
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

func (m *mongoRepository) FindCommute(ctx context.Context, date string) (*[]domain.Event, error) {
	var (
		recEvent []domain.Event
		err      error
	)

	t, _ := time.Parse(time.RFC3339, strings.Split(date, "T")[0]+"T00:00:00.000+00:00")

	opts := options.Find().SetSort(bson.D{{"start_date", 1}})

	cursor, err := m.Collection.Find(ctx, bson.D{{"start_date", bson.M{"$gte": primitive.NewDateTimeFromTime(t)}}, {"start_date", bson.M{"$lt": primitive.NewDateTimeFromTime(t.AddDate(0, 0, 1))}}}, opts)

	for cursor.Next(context.TODO()) {
		//Create a value into which the single document can beM decoded
		var elem domain.Event
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

func (m *mongoRepository) FindOneByID(ctx context.Context, id string) (*domain.Event, error) {
	var (
		recEvent domain.Event
		err      error
	)

	println("Searching | " + id)

	err = m.Collection.FindOne(ctx, bson.D{{"event_id", id}}).Decode(&recEvent)
	if err != nil {
		return &recEvent, err
	}

	return &recEvent, nil
}

func (m *mongoRepository) DeleteOneByID(ctx context.Context, id string) (*domain.Event, error) {
	var (
		recEvent domain.Event
		err      error
	)

	_, err = m.Collection.DeleteOne(ctx, bson.D{{"event_id", id}})
	if err != nil {
		return &recEvent, err
	}

	return &recEvent, nil
}

func (m *mongoRepository) UpdateOne(ctx context.Context, ev *domain.Event, eventID string) (*domain.Event, error) {

	var (
		event domain.Event
		err   error
	)

	u, err := m.Collection.UpdateOne(ctx, bson.D{{"event_id", eventID}, {"start_date", bson.M{"$gte": primitive.NewDateTimeFromTime(ev.StartDate)}}, {"start_date", bson.M{"$lt": primitive.NewDateTimeFromTime(ev.StartDate.AddDate(0, 0, 1))}}}, bson.D{{"$set", bson.D{{"location", ev.Location}, {"start_date", ev.StartDate}, {"end_date", ev.EndDate}, {"attendees", ev.Attendees}, {"mode", ev.Mode}, {"status", ev.Status}, {"type", ev.Type}, {"commute", ev.Commute}, {"remote_time", ev.RemoteTime}, {"remote_time_end", ev.RemoteTimeEnd}, {"summary", ev.Summary}}}})

	if u.MatchedCount <= 0 {
		_, err = m.Collection.InsertOne(ctx, bson.D{{"created_date", time.Now()}, {"event_id", ev.EventId}, {"start_date", ev.StartDate}, {"end_date", ev.EndDate}, {"mode", ev.Mode}, {"status", ev.Status}, {"type", ev.Type}, {"summary", ev.Summary}, {"attendees", ev.Attendees}, {"location", ev.Location}, {"commute", ev.Commute}, {"remote_time", ev.RemoteTime}, {"remote_time_end", ev.RemoteTimeEnd}, {"user_email", ev.UserEmail}})
	}

	if err != nil {
		println(err.Error())
		return &event, err
	}
	println("UPDATING...")

	return &event, nil
}

func (m *mongoRepository) DeleteAll(ctx context.Context, user_email string) (*[]domain.Event, error) {

	var (
		recEvent []domain.Event
		err      error
	)

	cursor, err := m.Collection.Find(ctx, bson.D{{"user_email", user_email}})

	for cursor.Next(context.TODO()) {
		//Create a value into which the single document can beM decoded
		var elem domain.Event
		err := cursor.Decode(&elem)
		if err != nil {
			println(err)
		}
		recEvent = append(recEvent, elem)
	}

	_, err = m.Collection.DeleteMany(ctx, bson.D{{"user_email", user_email}})
	if err != nil {
		println(err.Error())
		return &recEvent, err
	}
	println("UPDATING...")

	return &recEvent, nil
}
