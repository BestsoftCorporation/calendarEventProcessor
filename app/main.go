package main

import (
	"flag"
	_eventCalendar "github.com/bxcodec/go-clean-arch/event/delivery/calendar"
	_map "github.com/bxcodec/go-clean-arch/map/usecase"
	//"flag"
	"fmt"

	_eventRepo "github.com/bxcodec/go-clean-arch/event/repository/mongo"
	_eventUcase "github.com/bxcodec/go-clean-arch/event/usecase"
	_appRepo "github.com/bxcodec/go-clean-arch/integration/repository"
	_appUcase "github.com/bxcodec/go-clean-arch/integration/usecase"
	_recEventHttp "github.com/bxcodec/go-clean-arch/standard_trip/delivery/http"
	_recEventRepo "github.com/bxcodec/go-clean-arch/standard_trip/repository/mongo"
	recEventUsecase "github.com/bxcodec/go-clean-arch/standard_trip/usecase"
	echoSwagger "github.com/swaggo/echo-swagger"

	_tripRepo "github.com/bxcodec/go-clean-arch/trip/repository/mongo"
	_tripUcase "github.com/bxcodec/go-clean-arch/trip/usecase"

	//"google.golang.org/grpc"
	//"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"

	_userHttp "github.com/bxcodec/go-clean-arch/user/delivery/http"
	_userRepo "github.com/bxcodec/go-clean-arch/user/repository/mongo"
	_userUcase "github.com/bxcodec/go-clean-arch/user/usecase"

	_jwt "github.com/bxcodec/go-clean-arch/jwt/usecase"

	"github.com/bxcodec/go-clean-arch/bootstrap"
	/*
		_eventCalendar "github.com/bxcodec/go-clean-arch/event/delivery/calendar"
		_eventRepo "github.com/bxcodec/go-clean-arch/event/repository/mongo"
		_eventUcase "github.com/bxcodec/go-clean-arch/event/usecase"
		_tripRepo "github.com/bxcodec/go-clean-arch/trip/repository/mongo"
		_tripUcase "github.com/bxcodec/go-clean-arch/trip/usecase"
	*/
	_loginHttp "github.com/bxcodec/go-clean-arch/login/delivery/http"
	_loginUsecase "github.com/bxcodec/go-clean-arch/login/usecase"
	//pb "github.com/bxcodec/go-clean-arch/grpc"
)

func cors(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		return next(c)
	}
}

const (
	defaultName = "world"
)

//api.komon-beta.com
var (
	addr      = flag.String("addr", "localhost:50051", "the address to connect to")
	CompanyId = flag.Int("company_id", 1, "Name to greet")
	Id        = flag.Int("id", 1, "Name to greet")
)

// @title           Swagger Example API
func main() {
	e := echo.New()

	println("AWS TEST 90")

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	timeoutContext := time.Duration(bootstrap.App.Config.GetInt("context.timeout")) * time.Second

	database := bootstrap.App.Mongo.Database(bootstrap.App.Config.GetString("mongo.name"))

	e.Use(cors)

	appRepo := _appRepo.NewMongoAppRepository(*database)
	appUsecase := _appUcase.NewAppUsecase(appRepo, timeoutContext)

	userRepo := _userRepo.NewMongoRepository(*database)
	usrUsecase := _userUcase.NewUserUsecase(userRepo, timeoutContext)
	_userHttp.NewUserHandler(e, usrUsecase)

	recEveRepo := _recEventRepo.NewMongoRepository(*database)
	recEvUsecase := recEventUsecase.NewRecuringEventUsecase(recEveRepo, timeoutContext)
	_recEventHttp.NewStandardTripHandler(e, recEvUsecase, appUsecase)

	tripRepo := _tripRepo.NewMongoRepository(*database)
	tripUsecase := _tripUcase.NewTripUsecase(tripRepo, timeoutContext)

	eveRepo := _eventRepo.NewMongoRepository(*database)
	eveUsecase := _eventUcase.NewEventUsecase(eveRepo, timeoutContext)

	eveFetchRepo := _eventRepo.NewFetchEventMongoRepository(*database)
	eveFetchUsecase := _eventUcase.NewFetchEentUsecase(eveFetchRepo, timeoutContext)

	cacheRepo := _tripRepo.NewMongoCacheRepository(*database)
	cacheUsecase := _tripUcase.NewCacheUsecase(cacheRepo, timeoutContext)

	_map.NewMapsHandler(cacheUsecase)

	_eventCalendar.NewCalendarHandler(e, eveUsecase, tripUsecase, recEvUsecase, eveFetchUsecase)

	jwt := _jwt.NewJwtUsecase(userRepo, timeoutContext, bootstrap.App.Config)
	userJwt := e.Group("")
	jwt.SetJwtUser(userJwt)
	adminJwt := e.Group("")
	jwt.SetJwtUser(adminJwt)
	generalJwt := e.Group("")
	jwt.SetJwtUser(generalJwt)

	//Handle For login endpoint
	loginUsecase := _loginUsecase.NewLoginUsecase(userRepo, timeoutContext)
	_loginHttp.NewLoginHandler(e, loginUsecase, bootstrap.App.Config)

	appPort := fmt.Sprintf(":%s", bootstrap.App.Config.GetString("server.address"))
	log.Fatal(e.Start(appPort))
}
