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

type CompanyReportHandler struct {
	RecourEventUsecase domain.CompanyReportUsecase
}

func NewCompanyReportHandler(e *echo.Echo, uu domain.CompanyReportUsecase) {
	handler := &CompanyReportHandler{
		RecourEventUsecase: uu,
	}
	e.POST("/companyReport", handler.InsertOne)
	//e.GET("/companyReport", handler.FindOne)
}

func (CompanyReport *CompanyReportHandler) status(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
}

func (CompanyReport *CompanyReportHandler) InsertOne(c echo.Context) error {
	var (
		recEv domain.CompanyReport
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

	result, err := CompanyReport.RecourEventUsecase.InsertOne(ctx, &recEv)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

/*
func (CompanyReport *CompanyReportHandler) FindOne(c echo.Context) error {
	var (
		recEvent domain.Event
		err      error
	)


		err = m.Collection.FindOne(ctx, bson.M{"user_email": userEmail}).Decode(&recEvent)
		if err != nil {
			return &recEvent, err
		}

	return &recEvent, nil
}
*/

func isRequestValid(m *domain.CompanyReport) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}
