package usecase

import (
	"context"
	"github.com/bxcodec/go-clean-arch/domain"
	"googlemaps.github.io/maps"
	"log"
	"strconv"
	"strings"
	"time"
)

var handler = &MapsHandler{}

type MapsHandler struct {
	CacheUsecase domain.CacheUsecase
}

func NewMapsHandler(cu domain.CacheUsecase) {
	handler = &MapsHandler{
		CacheUsecase: cu,
	}

}

func GetLatLon(address string) (float64, float64) {
	c, err := maps.NewClient(maps.WithAPIKey("AIzaSyBYASTwYCLA_p6bXp35-dgbiuIy-ZBq7-A"))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	r := &maps.GeocodingRequest{
		Address: address,
	}

	var res []maps.GeocodingResult
	res, _ = c.Geocode(context.Background(), r)

	return res[0].Geometry.Location.Lat, res[0].Geometry.Location.Lng

}

func GetLocation(origin string, dest string, mode string, t time.Time, start int) int {
	c, err := maps.NewClient(maps.WithAPIKey("AIzaSyBYASTwYCLA_p6bXp35-dgbiuIy-ZBq7-A"))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	travelMode := maps.TravelModeTransit

	if mode == "WALKING" {
		travelMode = maps.TravelModeWalking
	} else if mode == "BICYCLING" {
		travelMode = maps.TravelModeBicycling
	} else if mode == "DRIVING" {
		travelMode = maps.TravelModeDriving
	}

	r := &maps.DirectionsRequest{}

	var currentTime time.Time

	if t.Before(time.Now()) {
		currentTime = time.Now().AddDate(0, 0, 7)
		day := strings.ToLower(t.Weekday().String()[:3])

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
	} else {
		currentTime = t
	}

	cache, _ := handler.CacheUsecase.FindOne(context.Background(), origin, dest, mode, strings.ToLower(t.Weekday().String()[:3]))

	if cache.StartLocation != "" && cache.Distance != 0 {
		return cache.Distance
	}

	if start == 1 {
		r = &maps.DirectionsRequest{
			Origin:       origin,
			Destination:  dest,
			ArrivalTime:  strconv.FormatInt(time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), currentTime.Hour(), currentTime.Minute(), 0, 0, time.UTC).Unix(), 10),
			Alternatives: false,
			Mode:         travelMode,
			Optimize:     true,
		}
	} else {
		r = &maps.DirectionsRequest{
			Origin:        origin,
			Destination:   dest,
			DepartureTime: strconv.FormatInt(time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), currentTime.Hour(), currentTime.Minute(), 0, 0, time.UTC).Unix(), 10),
			Alternatives:  false,
			Mode:          travelMode,
			Optimize:      true,
		}
	}

	route, _, err := c.Directions(context.Background(), r)
	if err != nil {
		println(err)
		log.Fatalf("fatal error: %s", err)
		return 0
	}

	sumDistance := 0
	if len(route) > 0 {
		steps := route[0].Legs[0].Steps
		for _, element := range steps {
			if element.TravelMode != "TRANSIT" {
				sumDistance += element.Distance.Meters
			}
		}

	}

	cache = &domain.Cache{
		StartLocation: origin,
		EndLocation:   dest,
		Distance:      sumDistance,
		Mode:          mode,
		Day:           strings.ToLower(t.Weekday().String()[:3]),
	}

	handler.CacheUsecase.InsertOne(context.Background(), cache)

	return sumDistance
}
