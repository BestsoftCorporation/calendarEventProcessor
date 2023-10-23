package usecase

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/bxcodec/go-clean-arch/domain"

	"github.com/google/uuid"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

func GetEvents(email string, minTime time.Time, deleted bool) *calendar.Events {

	ctx := context.Background()

	jsonCredentials, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Println(err)
	}

	config, err := google.JWTConfigFromJSON(jsonCredentials, "https://www.googleapis.com/auth/calendar")
	if err != nil {
		log.Println(err)
	}
	config.Subject = email

	ts := config.TokenSource(ctx)

	srv, err := calendar.NewService(ctx, option.WithTokenSource(ts))

	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	utc, err := time.LoadLocation("Europe/Brussels")
	fmt.Println("Location:", utc, ":Time:", time.Now().In(utc))
	events, err := srv.Events.List(email).ShowDeleted(deleted).
		SingleEvents(false).UpdatedMin(minTime.Add(+time.Second * 1).Format("2006-01-02T15:04:05Z")).OrderBy("updated").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}

	return events

}

func DeleteAllEvents(email string, fromEmail string) *calendar.Events {

	ctx := context.Background()

	jsonCredentials, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Println(err)
	}

	config, err := google.JWTConfigFromJSON(jsonCredentials, "https://www.googleapis.com/auth/calendar")
	if err != nil {
		log.Println(err)
	}
	config.Subject = "system1@komon.cloud"

	ts := config.TokenSource(ctx)

	srv, err := calendar.NewService(ctx, option.WithTokenSource(ts))

	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	fmt.Println("from email " + email)

	result := srv.Calendars.Clear("primary").Do()
	if true {
		println("Cleaning " + result.Error())
		return nil
	}
	events, err := srv.Events.List(email).MaxResults(1000).SingleEvents(false).ShowDeleted(false).Do()

	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}

	for _, ev := range events.Items {
		//println(ev.Summary)

		//if ev.Attendees[0].Email == fromEmail {
		//	println("deleting " + ev.Summary)
		DeleteEvent(email, ev.Id)
		println(ev.Summary)
		//}
	}

	println("removed")
	return events

}

func DeleteEvent(calendar_id string, event_id string) {

	ctx := context.Background()

	jsonCredentials, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Println(err)
	}

	config, err := google.JWTConfigFromJSON(jsonCredentials, "https://www.googleapis.com/auth/calendar")
	if err != nil {
		log.Println(err)
	}

	config.Subject = calendar_id

	ts := config.TokenSource(ctx)

	srv, err := calendar.NewService(ctx, option.WithTokenSource(ts))

	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	//srv.Calendars.Clear("system1@komon.cloud").Do()
	err = srv.Events.Delete(calendar_id, event_id).Do()
	if err != nil {
		println(err.Error())
	}
}

func UpdateEvent(calendar_id string, event_id string, event *calendar.Event) *calendar.Event {

	ctx := context.Background()

	jsonCredentials, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Println(err)
	}

	config, err := google.JWTConfigFromJSON(jsonCredentials, "https://www.googleapis.com/auth/calendar")
	if err != nil {
		log.Println(err)
	}

	config.Subject = calendar_id

	ts := config.TokenSource(ctx)

	srv, err := calendar.NewService(ctx, option.WithTokenSource(ts))

	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	_, err = srv.Events.Update(calendar_id, event_id, event).Do()
	if err != nil {
		println(err.Error())
		return nil
	}
	return event
}

func InsertEvent(calendar_id string, event *calendar.Event) {

	ctx := context.Background()

	jsonCredentials, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Println(err)
	}

	config, err := google.JWTConfigFromJSON(jsonCredentials, "https://www.googleapis.com/auth/calendar")
	if err != nil {
		log.Println(err)
	}

	config.Subject = calendar_id

	ts := config.TokenSource(ctx)

	srv, err := calendar.NewService(ctx, option.WithTokenSource(ts))

	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	_, err = srv.Events.Insert(calendar_id, event).Do()
	if err != nil {
		println(err.Error())
		return
	}
}

func PushOnboardingEvents(usecase domain.EventUsecase, email string, prefix string, day string, location string) time.Time {

	ctx := context.Background()

	mode := "🚌 > 🏢 > 🚌"

	if prefix == "car" {
		mode = "🚗  > 🏢 > 🚗"
	} else if prefix == "moto" {
		mode = "🏍️ > 🏢 > 🏍️"
	} else if prefix == "metro" || prefix == "tram" || prefix == "bus" || prefix == "train" {
		mode = "🚌  > 🏢 > 🚌️ "
	} else if prefix == "walk" || prefix == "step" {
		mode = "🚶 > 🏢 > 🚶️"
	} else if prefix == "bike" {
		mode = "🚲  > 🏢 > 🚲"
	} else if prefix == "remote" {
		mode = "🏠"
	} else if prefix == "off" {
		mode = "⛱ day off"
	}

	jsonCredentials, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Println(err)
	}

	config, err := google.JWTConfigFromJSON(jsonCredentials, "https://www.googleapis.com/auth/calendar")
	if err != nil {
		log.Println(err)
	}
	calendarId := "system1@komon.cloud" //prefix + "@komon.cloud" //prefix + "@komon.cloud"

	config.Subject = calendarId

	ts := config.TokenSource(ctx)

	srv, err := calendar.NewService(ctx, option.WithTokenSource(ts))

	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	currentTime := time.Now()

	println("Day ", int(currentTime.Weekday()))
	println(day+" ", int(time.Wednesday))
	println(int(time.Monday) - int(currentTime.Weekday()))

	if day == strings.ToLower(currentTime.Weekday().String()[:3]) {
		currentTime = currentTime.AddDate(0, 0, 0)
	} else if day == "mon" {
		currentTime = currentTime.AddDate(0, 0, int(time.Monday)-int(currentTime.Weekday()))
	} else if day == "tue" {
		currentTime = currentTime.AddDate(0, 0, int(time.Tuesday)-int(currentTime.Weekday()))
	} else if day == "wed" {
		currentTime = currentTime.AddDate(0, 0, int(time.Wednesday)-int(currentTime.Weekday()))
	} else if day == "thu" {
		currentTime = currentTime.AddDate(0, 0, int(time.Thursday)-int(currentTime.Weekday()))
	} else if day == "fri" {
		currentTime = currentTime.AddDate(0, 0, int(time.Friday)-int(currentTime.Weekday()))
	} else if day == "sat" {
		currentTime = currentTime.AddDate(0, 0, int(time.Saturday)-int(currentTime.Weekday()))
	} else if day == "sun" {
		currentTime = currentTime.AddDate(0, 0, int(time.Sunday)-int(currentTime.Weekday()))
	}

	currentTime = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 8, 0, 0, currentTime.Nanosecond(), currentTime.Location())

	att := []*calendar.EventAttendee{
		{Email: email},
	}

	if prefix == "remote" {
		att = []*calendar.EventAttendee{
			{Email: email},
			{Email: prefix + "@komon.cloud"},
		}
	}

	event := &calendar.Event{
		Attendees:   att,
		Description: "📍 Update your office address\n🚀 Add modal mail to change transport\n🏠 Add remote@komon.cloud to switch to remote\n🕓 Specify your remote time to combine with commute\n🛑 Decline to switch to day off\n",
		End: &calendar.EventDateTime{
			Date: currentTime.Format("2006-01-02"),
		},
		Start: &calendar.EventDateTime{
			Date: currentTime.Format("2006-01-02"),
		}, Status: "",
		Recurrence:      []string{"RRULE:FREQ=WEEKLY"},
		Location:        location,
		Summary:         mode,
		GuestsCanModify: true,
		Organizer:       &calendar.EventOrganizer{Email: calendarId},
		Creator:         &calendar.EventCreator{Email: calendarId},
		ServerResponse:  googleapi.ServerResponse{},
		Reminders:       &calendar.EventReminders{UseDefault: false},
		Transparency:    "transparent",
	}

	event, err = srv.Events.Insert(calendarId, event).SendUpdates("All").Do()
	if err != nil {
		log.Fatalf("Unable to create event. %v\n", err.Error())
	}

	ev := domain.Event{EventId: event.Id, Mode: prefix, StartDate: time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location()), EndDate: time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location()), Location: location, Summary: mode, UserEmail: email, Commute: true}

	usecase.UpdateOne(context.Background(), &ev, event.Id)

	timest, _ := time.Parse(time.RFC3339, event.Start.Date)
	ev = domain.Event{EventId: event.Id, StartDate: timest, Location: event.Location, Mode: prefix, Status: event.Status, Summary: event.Summary, Type: event.EventType, UserEmail: email, Commute: true}
	_, err = usecase.InsertOne(context.Background(), &ev)
	if err != nil {
		println(err.Error())
	}

	fmt.Printf("Event created: %s\n", event.HtmlLink)

	//uuid, _ := uuid.New()
	if prefix == "car1" {
		err := srv.Channels.Stop(&calendar.Channel{
			Id:         "01234567-89ab-cdef-0123456789ab56",
			ResourceId: "Dzbp8mr797i-gT_HE_Rqwd3jYaU",
		}).Do()
		if err != nil {
			println(err.Error())
			return time.Now()
		}
	}

	return currentTime

	//events, _ := srv.Events.Get(calendarId, "1d3c7351-e960-4942-9ce5-87fe9e19a62c").Do()

	//println(events.Summary)
	//println(channel.Address)

	//t := time.Now().Format(time.RFC3339)

}

func StartWatchEvents(email string) {
	ctx := context.Background()

	jsonCredentials, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Println(err)
	}

	config, err := google.JWTConfigFromJSON(jsonCredentials, "https://www.googleapis.com/auth/calendar")
	if err != nil {
		log.Println(err)
	}
	config.Subject = "jules@komon.co"

	ts := config.TokenSource(ctx)

	srv, err := calendar.NewService(ctx, option.WithTokenSource(ts))

	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	uu := uuid.New()

	eventsWatchCall := srv.Events.Watch(email, &calendar.Channel{
		Address:        "https://brain.komonqua.com/calendarUpdateNew",
		Expiration:     time.Now().AddDate(10, 0, 0).UnixNano() / int64(time.Millisecond),
		Id:             uu.String(),
		Kind:           "",
		Payload:        true,
		ResourceId:     "Dzbp8mr797i-gT_HE_Rqwd3jYaUY",
		ResourceUri:    "",
		Token:          "",
		Type:           "web_hook",
		ServerResponse: googleapi.ServerResponse{},
	})

	_, err = eventsWatchCall.Do()
	if err != nil {
		println(err.Error())
		//return nil
	}

	//println(channel.Expiration)
}
