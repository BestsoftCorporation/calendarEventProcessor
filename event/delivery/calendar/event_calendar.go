package calendar

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	_calendar "github.com/bxcodec/go-clean-arch/calendar/usecase"
	"github.com/bxcodec/go-clean-arch/domain"
	_map "github.com/bxcodec/go-clean-arch/map/usecase"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/googleapi"
)

type ResponseError struct {
	Message string `json:"message"`
}

var handler = &CalendarHandler{}

type CalendarHandler struct {
	Echo                *echo.Echo
	EventUsecase        domain.EventUsecase
	TripUsecase         domain.TripUsecase
	StandardTripUsecase domain.StandardTripUsecase
	EventFetchUsecase   domain.EventFetchUsecase
}

func PushEvents(email string, transport string, day string, location string) {

	dt := _calendar.PushOnboardingEvents(handler.EventUsecase, email, transport, day, location)

	if dt.Before(time.Now().AddDate(0, 0, 1)) {
		handler.generateTrips(dt, true, email)
	}

}

func NewCalendarHandler(e *echo.Echo, uu domain.EventUsecase, tr domain.TripUsecase, st domain.StandardTripUsecase, ef domain.EventFetchUsecase) {
	handler = &CalendarHandler{
		Echo:                e,
		EventUsecase:        uu,
		TripUsecase:         tr,
		StandardTripUsecase: st,
		EventFetchUsecase:   ef,
	}
	e.POST("/calendarUpdateNew", handler.canedarUpdates)
	e.GET("/fetchEvents", handler.fetchEvents)
	e.POST("/subscribe", handler.subscribe)
	e.DELETE("/deleteAllEvents", handler.deleteAllEvents)
	e.GET("/trips", handler.getTrips)
	//handler.processFutureEvents()
}

func (Calendar *CalendarHandler) canedarUpdates(c echo.Context) error {
	resource := c.Request().Header.Get("x-goog-resource-uri")
	prefix := strings.Split(resource, "/")
	pre := strings.Split(prefix[6], "%40")[0]
	if pre == "system2" || pre == "system1" {
		Calendar.GetEvents(pre)
	} else {
		Calendar.GetEventsPro(pre, calendar.Event{Summary: ""})
	}

	return c.JSON(http.StatusOK, "OK")
}

func (Calendar *CalendarHandler) fetchEvents(c echo.Context) error {

	Calendar.GetEvents("system1")
	//Calendar.GetEventsPro("car", calendar.Event{Summary: ""})
	return c.JSON(http.StatusOK, "OK")
}

func (Calendar *CalendarHandler) deleteAllEvents(c echo.Context) error {

	email := c.Request().Header.Get("email")
	Calendar.StandardTripUsecase.DeleteOne(context.Background(), email)
	events, _ := Calendar.EventUsecase.DeleteAll(context.Background(), email)
	Calendar.TripUsecase.DeleteAll(context.Background(), email)
	for _, ev := range *events {
		_calendar.DeleteEvent("system1@komon.cloud", ev.EventId)
	}
	//_calendar.DeleteAllEvents("system1@komon.cloud", email)

	return c.JSON(http.StatusOK, "OK")
}

func (Calendar *CalendarHandler) getTrips(c echo.Context) error {
	email := c.Request().Header.Get("email")
	trips, _ := Calendar.TripUsecase.FindAll(context.Background(), email)

	return c.JSON(http.StatusOK, trips)
}

func (Calendar *CalendarHandler) GetEventsPro(prefix string, eventForTomorrow calendar.Event) {

	eu := Calendar.EventUsecase
	tr := Calendar.TripUsecase
	recEvent := Calendar.StandardTripUsecase

	//utc, _ := time.LoadLocation("Europe/Brussels")
	var events calendar.Events

	if eventForTomorrow.Summary == "" {
		lf, _ := Calendar.EventFetchUsecase.FindOneByType(context.Background(), prefix)

		for lf.Type == "" {
			lf, _ = Calendar.EventFetchUsecase.FindOneByType(context.Background(), prefix)
		}

		println("difference", int64(time.Now().Sub(lf.LastFetch.Local()).Seconds()))

		if int64(time.Now().Sub(lf.LastFetch.Local()).Seconds()) < 2 {
			return
		}
		events.Items = _calendar.GetEvents(prefix+"@komon.cloud", lf.LastFetch, true).Items
		if len(events.Items) > 0 {
			fe := Calendar.EventFetchUsecase
			timest, _ := time.Parse(time.RFC3339, events.Items[len(events.Items)-1].Updated)
			fe.UpdateOne(context.Background(), &domain.EventFetch{LastFetch: timest}, lf.ID.Hex())
		}

	} else {
		events.Items = append(events.Items, &eventForTomorrow)
	}

	fmt.Println("Upcoming events:")
	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found.")
	} else {
		for _, item := range events.Items {

			if len(item.Attendees) == 0 || item.Creator.Email == "system1@komon.cloud" || item.Location == "" {
				continue
			}

			isUpdate, _ := eu.FindOneById(context.Background(), item.Id)
			tripPrevious, _ := tr.FindOne(context.Background(), item.Id, false)
			tripPreviousRet, _ := tr.FindOne(context.Background(), item.Id, true)

			addressUpdate := false
			modelUpdate := false
			if isUpdate != nil {
				if isUpdate.Location != item.Location {
					addressUpdate = true
				}
				if isUpdate.Mode != prefix {
					modelUpdate = true
				}
			}

			if item.Status == "cancelled" {
				tr.DeleteAllWhereLinkedId(context.Background(), item.Id)
				eu.DeleteOneByID(context.Background(), item.Id)
				continue
			}

			if item.Sequence != 0 || addressUpdate || modelUpdate {
				tr.DeleteAllWhereLinkedId(context.Background(), item.Id)
				eu.DeleteOneByID(context.Background(), item.Id)
				isUpdate = &domain.Event{}
				tripPrevious = &domain.Trip{}
				tripPreviousRet = &domain.Trip{}
			}

			println("Processing pro event: " + item.Creator.Email)
			date := item.Start.DateTime
			if date == "" {
				date = item.Start.Date

			}

			if isUpdate.Processed == true {
				continue
			}

			println("Previous trip " + tripPrevious.LinkedID)

			t, _ := time.Parse(time.RFC3339, item.Start.DateTime)
			tEnd, _ := time.Parse(time.RFC3339, item.End.DateTime)

			allAmount := 0
			mode := "WALKING"
			processed := true

			if prefix == "car" || prefix == "moto" {
				mode = "DRIVING"
			} else if prefix == "metro" || prefix == "tram" || prefix == "bus" || prefix == "train" {
				mode = "TRANSIT"
			} else if prefix == "walk" || prefix == "step" {
				mode = "WALKING"
			} else if prefix == "bike" {
				mode = "BICYCLING"
			}

			if t.Before(time.Now().AddDate(0, 0, 1)) {
				recEv, _ := recEvent.FindOne(context.Background(), item.Creator.Email)

				locPoint := recEv.UserWorkspaceAddress
				dest := locPoint

				for _, d := range recEv.Days {
					if d.Day == strings.ToLower(t.Weekday().String()[:3]) {
						if d.Mode == "remote" {
							locPoint = recEv.UserHomeAddress
							dest = recEv.UserHomeAddress
						}
						break
					}
				}

				prevEvents, _ := eu.Find(context.Background(), item.Creator.Email, item.Start.DateTime)

				startLoc := locPoint

				if startLoc == "" {
					continue
				}

				if len(*prevEvents) > 0 {
					old := ""

					var lastEvent domain.Event
					after := false
					for _, ev := range *prevEvents {

						if ev.EventId == item.Id {
							continue
						}

						difference := t.Sub(ev.EndDate)

						if ev.Commute == true {
							if ev.Mode == "remote" || ev.Mode == "off" {
								locPoint = recEv.UserHomeAddress
								dest = recEv.UserHomeAddress
							} else if strings.Contains(ev.Mode, "rem+") {
								if t.After(ev.RemoteTime) && t.Before(ev.RemoteTimeEnd) {
									locPoint = recEv.UserHomeAddress
								} else {
									locPoint = ev.Location
								}

								if tEnd.After(ev.RemoteTime) && tEnd.Before(ev.RemoteTimeEnd) {
									dest = recEv.UserHomeAddress
								} else {
									dest = ev.Location
								}

							} else {
								locPoint = ev.Location
								dest = ev.Location
							}
							continue
						}

						if ev.StartDate.After(tEnd) || ev.StartDate == tEnd {
							after = true
							difference = ev.StartDate.Sub(tEnd)
							startLoc = locPoint

							println("Time difference 2 ", int64(difference.Minutes()))

							if int64(difference.Minutes()) > 60 {
								dest = locPoint
							} else {
								dest = ev.Location
							}

							if old == "" {
								startLoc = locPoint
								break
							}

							difference = t.Sub(lastEvent.EndDate)
							if int64(difference.Minutes()) > 60 {
								startLoc = locPoint
							} else {
								startLoc = lastEvent.Location
							}

							break
						}

						after = false
						lastEvent = ev
						//oldTime = ev.EndDate
						old = ev.Location
					}

					if !after {
						difference := t.Sub(lastEvent.EndDate)
						if int64(difference.Minutes()) > 60 {
							startLoc = locPoint
							dest = locPoint
						} else {
							startLoc = lastEvent.Location
							dest = locPoint
						}
					}
				}

				latS, logS := _map.GetLatLon(startLoc)
				latD, logD := 0.0, 0.0
				if item.Location != "" {
					latD, logD = _map.GetLatLon(item.Location)
				}

				if tripPrevious.StartLocation != startLoc || tripPrevious.Destination != item.Location {
					//allAmount = _map.GetLocation(startLoc, item.Location, mode, t, 1)
				} else {
					//allAmount = tripPrevious.AllowanceAmount
				}

				for _, att := range item.Attendees {
					if strings.Split(att.Email, "@")[0] == "home" {
						startLoc = recEv.UserHomeAddress
						dest = recEv.UserHomeAddress
						break
					}
				}

				trip := domain.Trip{
					Date:            t.UTC(),
					StartTime:       t.UTC(),
					Disabled:        false,
					StartLocation:   startLoc,
					Destination:     item.Location,
					AllowanceAmount: allAmount,
					StartLocationGeo: domain.Geo{
						Lat: latS,
						Log: logS,
					},
					DestinationGeo: domain.Geo{
						Lat: latD,
						Log: logD,
					},
					LinkedID:  item.Id,
					Mode:      mode,
					RemoteDay: false,
					UserEmail: item.Creator.Email,
					Validate:  false,
					Returning: false,
					Commute:   false,
				}

				if tripPreviousRet.StartLocation != startLoc || tripPreviousRet.Destination != item.Location {
					//allAmount = _map.GetLocation(startLoc, item.Location, mode, tEnd, 0)
				} else {
					//allAmount = tripPreviousRet.AllowanceAmount
				}

				latDR, logDR := _map.GetLatLon(dest)

				returnTrip := domain.Trip{
					Date:            t.UTC(),
					StartTime:       tEnd.UTC(),
					Disabled:        false,
					StartLocation:   item.Location,
					Destination:     dest,
					AllowanceAmount: allAmount,
					StartLocationGeo: domain.Geo{
						Lat: latD,
						Log: logD,
					},
					DestinationGeo: domain.Geo{
						Lat: latDR,
						Log: logDR,
					},
					LinkedID:  item.Id,
					Mode:      mode,
					RemoteDay: false,
					UserEmail: item.Creator.Email,
					Validate:  false,
					Returning: true,
					Commute:   false,
				}

				if tripPrevious.UserEmail == "" && isUpdate.EventId == "" {
					tr.InsertOne(context.Background(), &trip)
					tr.InsertOne(context.Background(), &returnTrip)
				} else {
					tr.UpdateOne(context.Background(), &trip, tripPrevious.ID)
					tr.UpdateOne(context.Background(), &returnTrip, tripPreviousRet.ID)
				}
			} else {
				processed = false
			}

			event := domain.Event{
				EventId: item.Id, StartDate: t.UTC(),
				EndDate:   tEnd.UTC(),
				Commute:   false,
				Location:  item.Location,
				Mode:      mode,
				Status:    item.Status,
				Summary:   item.Summary,
				Type:      item.EventType,
				UserEmail: item.Creator.Email,
				Processed: processed,
			}

			if isUpdate.EventId == "" {
				eu.InsertOne(context.Background(), &event)
			} else {
				eu.UpdateOne(context.Background(), &event, event.EventId)
			}
		}
	}
}

func (Calendar *CalendarHandler) subscribe(c echo.Context) error {

	//_report.ExampleNewPDFGenerator()
	//_s3.S3Upload()

	_calendar.StartWatchEvents("system1@komon.cloud")

	_calendar.StartWatchEvents("car@komon.cloud")
	_calendar.StartWatchEvents("train@komon.cloud")
	//	_calendar.StartWatchEvents("bus@komon.cloud")
	_calendar.StartWatchEvents("remote@komon.cloud")
	_calendar.StartWatchEvents("moto@komon.cloud")
	_calendar.StartWatchEvents("bike@komon.cloud")
	_calendar.StartWatchEvents("walk@komon.cloud")
	return c.JSON(http.StatusOK, "OK")
}

func (CalendarHandler *CalendarHandler) processFutureEvents() {
	go func() {
		for {
			until, _ := CalendarHandler.EventFetchUsecase.FindOneByType(context.Background(), "COMMUTE")

			if until.Type == "" {
				break
			}

			if until.Processed == false {
				time.Sleep(time.Until(until.LastFetch))
			}

			users, _ := CalendarHandler.StandardTripUsecase.FindAll(context.Background())
			for _, user := range *users {
				CalendarHandler.generateTrips(until.LastFetch, false, user.UserEmail)
			}

			CalendarHandler.EventFetchUsecase.UpdateOne(context.Background(), &domain.EventFetch{
				LastFetch: until.LastFetch.AddDate(0, 0, 1),
				Type:      "PRO",
				Processed: false,
			}, until.ID.Hex())

		}
	}()

}

func (CalendarHandler *CalendarHandler) generateTrips(eventDate time.Time, now bool, email string) {

	var commute *[]domain.Event

	if now {
		// If its future day skip trip generating
		if eventDate.After(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()+1, 23, 59, 0, 0, eventDate.Location())) {
			return
		}
	}

	//Take all evens of user for that day
	commute, _ = CalendarHandler.EventUsecase.Find(context.Background(), email, eventDate.Format(time.RFC3339))
	hasCommute := false

	//Already has custom commute?
	for _, event := range *commute {
		if event.Commute == true {
			hasCommute = true
			break
		}
	}

	//If dont have commute add it
	if !hasCommute {
		user, _ := CalendarHandler.StandardTripUsecase.FindOne(context.Background(), email)

		for _, d := range user.Days {
			if d.Day == strings.ToLower(eventDate.Weekday().String()[:3]) {
				*commute = append(*commute, domain.Event{
					EndDate:    eventDate,
					StartDate:  eventDate,
					UpdateDate: eventDate,
					Location:   "",
					Mode:       d.Mode,
					UserEmail:  email,
					Processed:  false,
					Commute:    true,
				})
			}
		}
	}

	//Go through all events
	for _, event := range *commute {

		if event.Commute == true {

			//Take exsisting trips from db for this day
			tripToDeleteRet, _ := CalendarHandler.TripUsecase.FindOne(context.Background(), event.EventId, true)
			tripToDelete, _ := CalendarHandler.TripUsecase.FindOne(context.Background(), event.EventId, false)

			mode := "DRIVING"

			if event.Mode == "car" || event.Mode == "moto" {
				mode = "DRIVING"
			} else if event.Mode == "metro" || event.Mode == "tram" || event.Mode == "bus" || event.Mode == "train" {
				mode = "TRANSIT"
			} else if event.Mode == "walk" || event.Mode == "step" {
				mode = "WALKING"
			} else if event.Mode == "bike" {
				mode = "BICYCLING"
			}
			//Find this user by email
			user, _ := CalendarHandler.StandardTripUsecase.FindOne(context.Background(), event.UserEmail)
			location := event.Location
			//If user dont fill location its his default office
			if event.Location == "" {
				location = user.UserWorkspaceAddress
			}

			allAmount := 0
			remoteHours := 0

			println("Event mode " + event.Mode)

			if event.Mode != "remote" {
				if mode == "TRANSIT" && location != "" {
					allAmount = _map.GetLocation(user.UserHomeAddress, location, mode, event.StartDate, 1)
				}
			} else {

				remoteHours = 8
				mode = "remote"

				println("Remote hours ", remoteHours)
			}

			//Get lat and log for home
			latH, logH := _map.GetLatLon(user.UserHomeAddress)
			latD, logD := 0.0, 0.0
			if location != "" {
				latD, logD = _map.GetLatLon(location)
			} else {
				location = user.UserWorkspaceAddress
				latD, logD = _map.GetLatLon(user.UserWorkspaceAddress)
			}

			//If its a mixture of remote and commute
			if strings.Contains(event.Mode, "rem+") {
				remoteHours = int(event.RemoteTimeEnd.Sub(event.RemoteTime).Hours())

				if remoteHours == 0 {
					remoteHours = 1
				}

				oldRemoteTrip, _ := CalendarHandler.TripUsecase.FindOne(context.Background(), event.EventId+"_T", true)
				if len(oldRemoteTrip.BubbleId) > 0 {
					//CalendarHandler.EventUsecase.DeleteOneByID(context.Background(), event.EventId)
					CalendarHandler.TripUsecase.DeleteAllWhereLinkedId(context.Background(), event.EventId+"_T")
				}

				mode = strings.Split(event.Mode, "+")[1]

				remoteTrip := domain.Trip{
					Date:            event.StartDate,
					StartTime:       event.StartDate.Add(20 * time.Hour),
					Disabled:        false,
					StartLocation:   user.UserHomeAddress,
					Destination:     location,
					AllowanceAmount: allAmount,
					StartLocationGeo: domain.Geo{
						Lat: latH,
						Log: logH,
					},
					DestinationGeo: domain.Geo{
						Lat: latD,
						Log: logD,
					},
					LinkedID:    event.EventId + "_T",
					Mode:        "",
					RemoteDay:   false,
					UserEmail:   event.UserEmail,
					Validate:    false,
					Returning:   true,
					Commute:     true,
					RemoteHours: remoteHours,
				}
				CalendarHandler.TripUsecase.InsertOne(context.Background(), &remoteTrip)
				remoteHours = 0
			} else {
				remoteTrip, _ := CalendarHandler.TripUsecase.FindOne(context.Background(), event.EventId+"_T", true)
				if len(remoteTrip.BubbleId) > 0 {
					//CalendarHandler.EventUsecase.DeleteOneByID(context.Background(), event.EventId)
					CalendarHandler.TripUsecase.DeleteAllWhereLinkedId(context.Background(), event.EventId+"_T")
				}
			}

			trip := domain.Trip{
				Date:            event.StartDate,
				StartTime:       event.StartDate.Add(6 * time.Hour),
				Disabled:        false,
				StartLocation:   user.UserHomeAddress,
				Destination:     location,
				AllowanceAmount: allAmount,
				StartLocationGeo: domain.Geo{
					Lat: latH,
					Log: logH,
				},
				DestinationGeo: domain.Geo{
					Lat: latD,
					Log: logD,
				},
				LinkedID:    event.EventId,
				Mode:        mode,
				RemoteDay:   false,
				UserEmail:   event.UserEmail,
				Validate:    false,
				Returning:   false,
				Commute:     true,
				RemoteHours: remoteHours,
			}

			if event.Mode != "remote" && event.Mode != "off" {
				//allAmount = _map.GetLocation(user.UserWorkspaceAddress, user.UserHomeAddress, mode, event.EndDate, 0)

				//Returning trip
				returnTrip := domain.Trip{
					Date:          event.EndDate,
					StartTime:     event.EndDate,
					Disabled:      false,
					StartLocation: user.UserWorkspaceAddress,
					Destination:   location,
					StartLocationGeo: domain.Geo{
						Lat: latD,
						Log: logD,
					},
					DestinationGeo: domain.Geo{
						Lat: latH,
						Log: logH,
					},
					AllowanceAmount: allAmount,
					LinkedID:        event.EventId,
					Mode:            mode,
					RemoteDay:       false,
					UserEmail:       event.UserEmail,
					Validate:        false,
					Returning:       true,
					Commute:         true,
					RemoteHours:     remoteHours,
				}

				if len(tripToDeleteRet.LinkedID) > 0 {
					CalendarHandler.TripUsecase.UpdateOne(
						context.Background(),
						&returnTrip,
						tripToDeleteRet.ID)
				} else {
					CalendarHandler.TripUsecase.InsertOne(context.Background(), &returnTrip)
				}
			} else {
				if len(tripToDeleteRet.BubbleId) > 0 {

					tripToDeleteRet, _ = CalendarHandler.TripUsecase.DeleteByReturningAndID(
						context.Background(),
						true,
						event.EventId,
						"returning",
						event.EventId,
						true)

					tripToDelete, _ = CalendarHandler.TripUsecase.DeleteByReturningAndID(
						context.Background(),
						false,
						event.EventId,
						"returning",
						event.EventId,
						true)
				}
			}

			if event.Mode != "off" {
				if len(tripToDelete.LinkedID) > 0 {
					CalendarHandler.TripUsecase.UpdateOne(context.Background(), &trip, tripToDelete.ID)
				} else {
					CalendarHandler.TripUsecase.InsertOne(context.Background(), &trip)
				}
			} else {
				remoteTrip, _ := CalendarHandler.TripUsecase.FindOne(context.Background(), event.EventId+"_T", true)
				if len(remoteTrip.BubbleId) > 0 {
					CalendarHandler.TripUsecase.DeleteAllWhereLinkedId(context.Background(), event.EventId+"_T")
				}
				if len(tripToDelete.BubbleId) > 0 {
					CalendarHandler.TripUsecase.DeleteByReturningAndID(context.Background(),
						false, event.EventId,
						"returning",
						event.EventId,
						true)
					CalendarHandler.TripUsecase.DeleteByReturningAndID(context.Background(),
						true, event.EventId,
						"returning",
						event.EventId,
						true)
				}

			}

			println("Inserting trip")

		} else {
			CalendarHandler.TripUsecase.DeleteAllWhereLinkedId(context.Background(), event.EventId)

			if event.Processed == false {
				CalendarHandler.GetEventsPro(event.Mode, calendar.Event{
					Creator:   &calendar.EventCreator{Email: event.UserEmail},
					End:       &calendar.EventDateTime{DateTime: event.EndDate.String()},
					Id:        event.EventId,
					Location:  event.Location,
					Organizer: &calendar.EventOrganizer{Email: event.UserEmail},
					Start:     &calendar.EventDateTime{DateTime: event.StartDate.String()},
					Summary:   event.Summary,
				})
			}

		}

	}

}

func (Calendar *CalendarHandler) GetEvents(prefix string) {

	eu := Calendar.EventUsecase
	tr := Calendar.TripUsecase
	recEvent := Calendar.StandardTripUsecase
	fe := Calendar.EventFetchUsecase

	lf, _ := Calendar.EventFetchUsecase.FindOneByType(context.Background(), "DLT")

	for lf.Type == "" {
		lf, _ = Calendar.EventFetchUsecase.FindOneByType(context.Background(), "DLT")
	}
	/*
		if lf.Type == "" {
			Calendar.EventFetchUsecase.InsertOne(context.Background(), &domain.EventFetch{
				LastFetch: time.Now(),
				Type:      "DLT",
				Processed: false,
			})
		}

	*/

	println("difference", int64(time.Now().Sub(lf.LastFetch.Local()).Seconds()))

	if int64(time.Now().Sub(lf.LastFetch.Local()).Seconds()) < 2 {
		return
	}
	events := _calendar.GetEvents(prefix+"@komon.cloud", lf.LastFetch, false)

	println("last update ", lf.LastFetch.Local().String())
	println("time now ", time.Now().String())

	fmt.Println("Upcoming events:")
	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found.")
	} else {
		timest, _ := time.Parse(time.RFC3339, events.Items[len(events.Items)-1].Updated)

		var wg sync.WaitGroup

		for _, item := range events.Items {

			if len(item.Attendees) == 0 || item.Creator.Email != "system1@komon.cloud" {
				continue
			}
			println(item.Summary + " " + item.Attendees[0].ResponseStatus)
			user, _ := recEvent.FindOne(context.Background(), item.Attendees[0].Email)

			wg.Add(1)

			go func() {
				defer wg.Done()
				Calendar.processCommuteEvent(eu, tr, item, user)
			}()
		}

		wg.Wait()
		fe.UpdateOne(context.Background(), &domain.EventFetch{LastFetch: timest}, lf.ID.Hex())
	}

}

func (Calendar *CalendarHandler) processCommuteEvent(eu domain.EventUsecase, tr domain.TripUsecase, item *calendar.Event, user *domain.StandardTrip) {

	if item.Recurrence != nil && len(item.Recurrence) > 0 {
		if strings.Contains(item.Recurrence[0], "INTERVAL=2") && !strings.Contains(item.Description, "*") {
			/*
				tm, _ := time.Parse("2006-01-02", item.Start.Date)

				mode := "🚗"
				location := user.UserWorkspaceAddress
				for _, d := range user.Days {
					if d.Day == strings.ToLower(tm.Weekday().String()[:3]) {
						mode = d.Mode
						break
					}
				}
				timest, _ := time.Parse("2006-01-02", item.Start.Date)

				att := []*calendar.EventAttendee{
					{Email: item.Attendees[0].Email},
				}

				if mode == "remote" {
					att = []*calendar.EventAttendee{
						{Email: item.Attendees[0].Email},
						{Email: "remote@komon.cloud"},
					}
					location = user.UserHomeAddress

				}

					_calendar.InsertEvent(item.Creator.Email, &calendar.Event{
						Attendees:   att,
						Description: "* ✏️ Modify to change mode or workplace address\n🕓 Specify your telework time to combine with commute\n✅ Accept to confirm\n❌ Decline to switch to teleworking\n🗑 Delete to switch to day off\n *",
						End: &calendar.EventDateTime{
							Date:     timest.Add(time.Hour * 24 * 7).Format("2006-01-02"),
							TimeZone: "Europe/Brussels",
						},
						Start: &calendar.EventDateTime{
							Date:     timest.Add(time.Hour * 24 * 7).Format("2006-01-02"),
							TimeZone: "Europe/Brussels",
						},
						Location:        location,
						Recurrence:      item.Recurrence,
						Status:          "",
						Summary:         getIcons(mode),
						GuestsCanModify: true,
						Organizer:       &calendar.EventOrganizer{Email: item.Organizer.Email},
						Creator:         &calendar.EventCreator{Email: item.Organizer.Email},
						ServerResponse:  googleapi.ServerResponse{},
						EventType:       item.EventType,
					})
			*/
			println(item.Recurrence[0])
		}
	}

	if item.Start.DateTime != "" && strings.Contains(item.Creator.Email, "@komon.cloud") {

		timest, _ := time.Parse(time.RFC3339, item.Start.DateTime)
		timestEnd, _ := time.Parse(time.RFC3339, item.End.DateTime)

		desc := "📍 Update your office address\n🚀 Add modal mail to change transport\n❌ Remove remote@komon.cloud to switch to commute\n🛑 Decline to switch to day off\n"

		location := user.UserWorkspaceAddress
		if item.Location != "" {
			location = item.Location
		}

		if item.Summary == "⛱ day off" {
			location = ""
			desc = "❌ Remove remote@komon.cloud to switch to commute\n🛑 Decline to switch to day off\n"
		}

		var arr []string
		for _, att := range item.Attendees {
			arr = append(arr, att.Email)
		}

		startT, _ := time.Parse(time.RFC3339, strings.Split(timest.Format(time.RFC3339Nano), "T")[0]+"T00:00:00.000+00:00")

		eventPrev, _ := eu.FindOneById(context.Background(), strings.Split(item.Id, "_")[0])

		if "rem+"+eventPrev.Mode == "rem+remote" {
			eu.DeleteOneByID(context.Background(), eventPrev.EventId)
			return
		}

		if strings.Contains(item.Summary, "🏠 + ") {
			item.Summary = strings.Replace(item.Summary, "🏠 + ", "", 3)
			eu.DeleteOneByID(context.Background(), eventPrev.EventId)
		}

		ev := domain.Event{EventId: strings.Split(item.Id, "_")[0], Mode: "rem+" + eventPrev.Mode, Attendees: arr, StartDate: startT, EndDate: startT, Location: location, Summary: "🏠 + " + item.Summary, UserEmail: item.Attendees[0].Email, Commute: true, RemoteTime: timest.UTC(), RemoteTimeEnd: timestEnd.UTC()}

		_, err := eu.UpdateOne(context.Background(), &ev, strings.Split(item.Id, "_")[0])

		_calendar.UpdateEvent(item.Creator.Email, item.Id, &calendar.Event{
			Attendees: []*calendar.EventAttendee{
				{Email: item.Attendees[0].Email},
				{Email: "remote@komon.cloud"},
			},
			Description: desc,
			End: &calendar.EventDateTime{
				Date: timest.AddDate(0, 0, 1).Format("2006-01-02"),
			},
			Start: &calendar.EventDateTime{
				Date: timest.Format("2006-01-02"),
			}, Status:       "confirmed",
			Recurrence:      item.Recurrence,
			Location:        location,
			Summary:         "🏠 + " + item.Summary,
			GuestsCanModify: true,
			Organizer:       &calendar.EventOrganizer{Email: item.Organizer.Email},
			Creator:         &calendar.EventCreator{Email: item.Organizer.Email},
			ServerResponse:  googleapi.ServerResponse{},
			Reminders:       &calendar.EventReminders{UseDefault: false},
			Transparency:    "transparent",
		})

		Calendar.generateTrips(timest, true, item.Attendees[0].Email)

		if err != nil {
			println(err.Error())
		}

		return
	}

	if item.Attendees[0].ResponseStatus == "declined" && item.Summary != "⛱ day off" {
		timest, _ := time.Parse("2006-01-02", item.Start.Date)

		println("this declined")
		var arr []string
		for _, att := range item.Attendees {
			arr = append(arr, att.Email)
		}
		ev := domain.Event{EventId: strings.Split(item.Id, "_")[0], Location: "", Mode: "off", Attendees: arr, StartDate: timest, EndDate: timest, Summary: "⛱ day off", UserEmail: item.Attendees[0].Email, Commute: true}
		_, err := eu.UpdateOne(context.Background(), &ev, strings.Split(item.Id, "_")[0])

		_calendar.UpdateEvent(item.Creator.Email, item.Id, &calendar.Event{
			Attendees: []*calendar.EventAttendee{
				{Email: item.Attendees[0].Email},
			},
			Description: "🏠 Add remote@komon.cloud to switch to remote\n🕓 Specify your remote time to combine with day off\n🛑 Decline to switch to commute\n",
			End: &calendar.EventDateTime{
				Date: timest.AddDate(0, 0, 1).Format("2006-01-02"),
			},
			Start: &calendar.EventDateTime{
				Date: timest.Format("2006-01-02"),
			},
			Recurrence:      item.Recurrence,
			Status:          "confirmed",
			Summary:         "⛱ day off",
			GuestsCanModify: true,
			Organizer:       &calendar.EventOrganizer{Email: item.Organizer.Email},
			Creator:         &calendar.EventCreator{Email: item.Organizer.Email},
			ServerResponse:  googleapi.ServerResponse{},
			Reminders:       &calendar.EventReminders{UseDefault: false},
			Transparency:    "transparent",
		})

		Calendar.generateTrips(timest, true, item.Attendees[0].Email)
		if err != nil {
			println(err.Error())
		}
		return
	} else if item.Attendees[0].ResponseStatus == "declined" && item.Summary == "⛱ day off" {

		tm, _ := time.Parse("2006-01-02", item.Start.Date)
		mode := "car"
		location := user.UserWorkspaceAddress

		for _, d := range user.Days {
			if d.Day == strings.ToLower(tm.Weekday().String()[:3]) {
				mode = d.Mode
				if mode == "remote" {
					location = ""
				}
				break
			}
		}

		timest, _ := time.Parse("2006-01-02", item.Start.Date)

		println("this declined")
		var arr []string
		for _, att := range item.Attendees {
			arr = append(arr, att.Email)
		}

		if mode == "rem+remote" {
			return
		}

		ev := domain.Event{EventId: strings.Split(item.Id, "_")[0], Mode: mode, Attendees: arr, StartDate: timest, EndDate: timest, Location: location, Summary: getIcons(mode), UserEmail: item.Attendees[0].Email, Commute: true}
		_, err := eu.UpdateOne(context.Background(), &ev, strings.Split(item.Id, "_")[0])

		_calendar.UpdateEvent(item.Creator.Email, item.Id, &calendar.Event{
			Attendees: []*calendar.EventAttendee{
				{Email: item.Attendees[0].Email},
			},
			Description: "📍 Update your office address\n🚀 Add modal mail to change transport\n🏠 Add remote@komon.cloud to switch to remote\n🕓 Specify your remote time to combine with commute\n🛑 Decline to switch to day off\n",
			End: &calendar.EventDateTime{
				Date: timest.AddDate(0, 0, 1).Format("2006-01-02"),
			},
			Start: &calendar.EventDateTime{
				Date: timest.Format("2006-01-02"),
			},
			Location:        location,
			Recurrence:      item.Recurrence,
			Status:          "confirmed",
			Summary:         getIcons(mode),
			GuestsCanModify: true,
			Organizer:       &calendar.EventOrganizer{Email: item.Organizer.Email},
			Creator:         &calendar.EventCreator{Email: item.Organizer.Email},
			ServerResponse:  googleapi.ServerResponse{},
			Reminders:       &calendar.EventReminders{UseDefault: false},
			Transparency:    "transparent",
		})

		Calendar.generateTrips(timest, true, item.Attendees[0].Email)
		if err != nil {
			println(err.Error())
		}

		return
	}

	if len(item.Attendees) == 1 && strings.Contains(item.Summary, "🏠") {
		timest, _ := time.Parse("2006-01-02", item.Start.Date)
		location := ""

		newSummary := ""
		mode := "car"
		if strings.Contains(item.Summary, "🏠 + ") && !strings.Contains(item.Summary, "⛱ day off") {
			newSummary = strings.Replace(item.Summary, "🏠 + ", "", 3)
			if item.Location != "" {
				location = item.Location
			}
			eventPrev, _ := eu.FindOneById(context.Background(), strings.Split(item.Id, "_")[0])
			mode = strings.Replace(eventPrev.Mode, "rem+", "", 4)

		} else {
			for _, d := range user.Days {
				if d.Day == strings.ToLower(timest.Weekday().String()[:3]) {
					newSummary = d.Mode
					if newSummary == "remote" {
						location = ""
					} else {
						location = user.UserWorkspaceAddress
					}
					break
				}
			}
			newSummary = getIcons(newSummary)
		}

		var arr []string
		for _, att := range item.Attendees {
			arr = append(arr, att.Email)
		}

		if mode == "rem+remote" {
			return
		}

		ev := domain.Event{EventId: strings.Split(item.Id, "_")[0], Mode: mode, StartDate: timest, Attendees: arr, EndDate: timest, Location: location, Summary: newSummary, UserEmail: item.Attendees[0].Email, Commute: true}

		_, err := eu.UpdateOne(context.Background(), &ev, strings.Split(item.Id, "_")[0])

		println("Removing house " + newSummary)
		_calendar.UpdateEvent(item.Creator.Email, item.Id, &calendar.Event{
			Attendees: []*calendar.EventAttendee{
				{Email: item.Attendees[0].Email},
			},
			Description: "📍 Update your office address\n🚀 Add modal mail to change transport\n🏠 Add remote@komon.cloud to switch to remote\n🕓 Specify your remote time to combine with commute\n🛑 Decline to switch to day off\n",
			End: &calendar.EventDateTime{
				Date: timest.AddDate(0, 0, 1).Format("2006-01-02"),
			},
			Start: &calendar.EventDateTime{
				Date: timest.Format("2006-01-02"),
			},
			Location:        location,
			Status:          "confirmed",
			Recurrence:      item.Recurrence,
			Summary:         newSummary,
			GuestsCanModify: true,
			Organizer:       &calendar.EventOrganizer{Email: item.Organizer.Email},
			Creator:         &calendar.EventCreator{Email: item.Organizer.Email},
			ServerResponse:  googleapi.ServerResponse{},
			Reminders:       &calendar.EventReminders{UseDefault: false},
			Transparency:    "transparent",
		})

		Calendar.generateTrips(timest, true, item.Attendees[0].Email)
		if err != nil {
			println(err.Error())
		}

		return
	}

	if len(item.Attendees) == 2 {

		location := user.UserWorkspaceAddress
		if item.Location != "" {
			location = item.Location
		}

		desc := "📍 Update your office address\n🚀 Add modal mail to change transport\n🏠 Add remote@komon.cloud to switch to remote\n🕓 Specify your remote time to combine with commute\n🛑 Decline to switch to day off\n"

		att := []*calendar.EventAttendee{
			{Email: item.Attendees[0].Email},
		}
		if strings.Split(item.Attendees[1].Email, "@")[0] == "remote" && !strings.Contains(item.Summary, "🏠") {
			att = []*calendar.EventAttendee{
				{Email: item.Attendees[0].Email},
				{Email: "remote@komon.cloud"},
			}
			location = ""
			desc = "❌ Remove remote@komon.cloud to switch to commute\n🛑 Decline to switch to day off\n"
		} else if strings.Split(item.Attendees[1].Email, "@")[0] == "remote" {
			return
		}

		timest, _ := time.Parse("2006-01-02", item.Start.Date)

		var arr []string
		for _, att := range item.Attendees {
			arr = append(arr, att.Email)
		}
		ev := domain.Event{EventId: strings.Split(item.Id, "_")[0], Attendees: arr, Mode: strings.Split(item.Attendees[1].Email, "@")[0], StartDate: timest, EndDate: timest, Location: location, Summary: getIcons(strings.Split(item.Attendees[1].Email, "@")[0]), UserEmail: item.Attendees[0].Email, Commute: true}

		_, err := eu.UpdateOne(context.Background(), &ev, strings.Split(item.Id, "_")[0])

		Calendar.generateTrips(timest, true, item.Attendees[0].Email)

		_calendar.UpdateEvent(item.Creator.Email, item.Id, &calendar.Event{
			Attendees:   att,
			Description: desc,
			End: &calendar.EventDateTime{
				Date: timest.AddDate(0, 0, 1).Format("2006-01-02"),
			},
			Start: &calendar.EventDateTime{
				Date: timest.Format("2006-01-02"),
			}, Status:       "confirmed",
			Recurrence:      item.Recurrence,
			Location:        location,
			Summary:         getIcons(strings.Split(item.Attendees[1].Email, "@")[0]),
			GuestsCanModify: true,
			Organizer:       &calendar.EventOrganizer{Email: item.Organizer.Email},
			Creator:         &calendar.EventCreator{Email: item.Organizer.Email},
			ServerResponse:  googleapi.ServerResponse{},
			Reminders:       &calendar.EventReminders{UseDefault: false},
			Transparency:    "transparent",
		})

		if err != nil {
			println(err.Error())
		}

		return

	} else if len(item.Attendees) == 3 {

		firstModel := getSingleIcons(strings.Split(item.Attendees[1].Email, "@")[0])
		secondModel := getSingleIcons(strings.Split(item.Attendees[2].Email, "@")[0])

		summary := firstModel + " > 🏢 > " + secondModel
		att := []*calendar.EventAttendee{
			{Email: item.Attendees[0].Email},
		}

		desc := "📍 Update your office address\n🚀 Add modal mail to change transport\n🏠 Add remote@komon.cloud to switch to remote\n🕓 Specify your remote time to combine with commute\n🛑 Decline to switch to day off\n"

		location := user.UserWorkspaceAddress

		if item.Location != "" {
			location = item.Location
		}

		mode := strings.Split(item.Attendees[2].Email, "@")[0]
		eventPrev, _ := eu.FindOneById(context.Background(), strings.Split(item.Id, "_")[0])

		if firstModel == "🏠" {
			mode = "rem+" + strings.Split(item.Attendees[2].Email, "@")[0]
			summary = "🏠 + " + secondModel + " > 🏢 > " + secondModel
			att = []*calendar.EventAttendee{
				{Email: item.Attendees[0].Email},
				{Email: "remote@komon.cloud"},
			}

			desc = "📍 Update your office address\n🚀 Add modal mail to change transport\n❌ Remove remote@komon.cloud to switch to commute\n🛑 Decline to switch to day off\n"

		} else if secondModel == "🏠" {
			mode = "rem+" + strings.Split(item.Attendees[1].Email, "@")[0]
			summary = "🏠 + " + firstModel + " > 🏢 > " + firstModel
			att = []*calendar.EventAttendee{
				{Email: item.Attendees[0].Email},
				{Email: "remote@komon.cloud"},
			}

			desc = "📍 Update your office address\n🚀 Add modal mail to change transport\n❌ Remove remote@komon.cloud to switch to commute\n🛑 Decline to switch to day off\n"
		}

		timest, _ := time.Parse("2006-01-02", item.Start.Date)
		var arr []string
		for _, att := range item.Attendees {
			arr = append(arr, att.Email)
		}

		if mode == "rem+remote" {
			return
		}

		ev := domain.Event{EventId: strings.Split(item.Id, "_")[0], Attendees: arr, Mode: mode, StartDate: timest,
			EndDate: timest, Location: location, Summary: summary, UserEmail: item.Attendees[0].Email, Commute: true,
			RemoteTime: eventPrev.RemoteTime, RemoteTimeEnd: eventPrev.RemoteTimeEnd}

		_, err := eu.UpdateOne(context.Background(), &ev, strings.Split(item.Id, "_")[0])

		Calendar.generateTrips(timest, true, item.Attendees[0].Email)

		_calendar.UpdateEvent(item.Creator.Email, item.Id, &calendar.Event{
			Attendees:   att,
			Description: desc,
			End: &calendar.EventDateTime{
				Date: timest.AddDate(0, 0, 1).Format("2006-01-02"),
			},
			Start: &calendar.EventDateTime{
				Date: timest.Format("2006-01-02"),
			}, Status:       "confirmed",
			Recurrence:      item.Recurrence,
			Location:        location,
			Summary:         summary,
			GuestsCanModify: true,
			Organizer:       &calendar.EventOrganizer{Email: item.Organizer.Email},
			Creator:         &calendar.EventCreator{Email: item.Organizer.Email},
			ServerResponse:  googleapi.ServerResponse{},
			Reminders:       &calendar.EventReminders{UseDefault: false},
			Transparency:    "transparent",
		})

		if err != nil {
			println(err.Error())
		}

		println("updating")

		return
	}

	eventPrev, _ := eu.FindOneById(context.Background(), strings.Split(item.Id, "_")[0])

	if eventPrev.Location != item.Location {
		ev := domain.Event{EventId: strings.Split(item.Id, "_")[0], Attendees: eventPrev.Attendees, Mode: eventPrev.Mode, StartDate: eventPrev.StartDate, EndDate: eventPrev.EndDate, Location: item.Location, Summary: eventPrev.Summary, UserEmail: item.Attendees[0].Email, Commute: true}

		_, err := eu.UpdateOne(context.Background(), &ev, strings.Split(item.Id, "_")[0])
		timest, _ := time.Parse("2006-01-02", item.Start.Date)

		Calendar.generateTrips(timest, true, item.Attendees[0].Email)
		if err != nil {
			println(err.Error())
		}
	}

	/*
		if len(item.Attendees) == 1 {
			location := user.UserWorkspaceAddress
			if item.Location != "" {
				location = item.Location
			}

			prevEvent, err := eu.FindOneById(context.Background(), item.Id)
			prevEvent.Location = location
			_, err = eu.UpdateOne(context.Background(), prevEvent, strings.Split(item.Id, "_")[0])
			Calendar.generateTrips(prevEvent.StartDate, true)
			println("Updating address")
			if err != nil {
				println(err.Error())
			}
		}
	*/

}

func getIcons(prefix string) string {
	mode := "🚌 > 🏢 > 🚌"

	if prefix == "car" {
		mode = "🚗  > 🏢 > 🚗"
	} else if prefix == "moto" {
		mode = "🛵 > 🏢 > 🛵"
	} else if prefix == "metro" || prefix == "tram" || prefix == "bus" {
		mode = "🚌  > 🏢 > 🚌️ "
	} else if prefix == "train" {
		mode = "🚂 > 🏢 > 🚂️"
	} else if prefix == "walk" {
		mode = "🚶 > 🏢 > 🚶️"
	} else if prefix == "step" {
		mode = "🛴 > 🏢 > 🛴"
	} else if prefix == "bike" {
		mode = "🚲  > 🏢 > 🚲"
	} else if prefix == "remote" {
		mode = "🏠"
	}

	return mode
}

func getSingleIcons(prefix string) string {
	mode := "🚌"

	if prefix == "car" {
		mode = "🚗"
	} else if prefix == "moto" {
		mode = "🛵️"
	} else if prefix == "metro" || prefix == "tram" || prefix == "bus" {
		mode = "🚌"
	} else if prefix == "walk" {
		mode = "🚶"
	} else if prefix == "bike" {
		mode = "🚲"
	} else if prefix == "remote" {
		mode = "🏠"
	} else if prefix == "train" {
		mode = "🚂"
	} else if prefix == "step" {
		mode = "🛴"
	}

	return mode
}
