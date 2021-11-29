// Package request contains everything for request and response. Handler uses it
package request

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"net/http"
)

// CustomValidator replaces default validator of Echo
type CustomValidator struct {
	Validator *validator.Validate
}

// Validate defines the validity of struct
func (c *CustomValidator) Validate(i interface{}) error {
	if err := c.Validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

// CatID RequestCatID gets id
type CatID struct {
	ID int32 `param:"id" validate:"required,numeric,gt=0"`
}
