package request

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (c *CustomValidator) Validate(i interface{}) error {
	if err := c.Validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

type RequestCatId struct {
	ID int32 `param:"id" json:"id" bson:"id" query:"id" header:"id" form:"id" xml:"id" validate:"required,numeric,gt=0"`
}
