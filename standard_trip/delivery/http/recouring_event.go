package http

import (
	"context"
	"github.com/bxcodec/go-clean-arch/domain"
	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
)

type ResponseError struct {
	Message string `json:"message"`
}

type StandardTripHandler struct {
	RecourEventUsecase domain.StandardTripUsecase
	AppUsecase         domain.AppUsecase
}

func NewStandardTripHandler(e *echo.Echo, uu domain.StandardTripUsecase, app domain.AppUsecase) {
	handler := &StandardTripHandler{
		RecourEventUsecase: uu,
		AppUsecase:         app,
	}
	e.POST("/app", handler.InsertApp)
	e.POST("/calendarUpdate", handler.canedarUpdates)
	e.POST("/user", handler.InsertOne)
	e.GET("/status", handler.status)
}

func (StandardTrip *StandardTripHandler) canedarUpdates(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
}

func (StandardTrip *StandardTripHandler) eventsWatch(c echo.Context) error {

	return c.JSON(http.StatusOK, "OK")
}

func (StandardTrip *StandardTripHandler) status(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
}

func (StandardTrip *StandardTripHandler) InsertApp(c echo.Context) error {

	token := c.Request().Header.Get("Authorization")

	if token != "iujnxygfnashdg213ascxisax213" {
		return c.JSON(http.StatusUnauthorized, "Your app is not authorized to send requests!")
	}

	var (
		recEv domain.App
		err   error
	)

	err = c.Bind(&recEv)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := StandardTrip.AppUsecase.InsertOne(ctx, &recEv)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

func (StandardTrip *StandardTripHandler) InsertOne(c echo.Context) error {

	token := c.Request().Header.Get("Authorization")
	app, _ := StandardTrip.AppUsecase.FindOne(context.Background(), token)
	if app.AppName != "Bubble" {
		return c.JSON(http.StatusUnauthorized, "Your app is not authorized to send requests!")
	}

	var (
		recEv domain.StandardTrip
		err   error
	)

	err = c.Bind(&recEv)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	var ok bool
	if ok, err = isRequestValid(&recEv); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := StandardTrip.RecourEventUsecase.InsertOne(ctx, &recEv)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

func isRequestValid(m *domain.StandardTrip) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}
