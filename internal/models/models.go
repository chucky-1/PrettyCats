package models

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Cats struct {
	ID   int32  `json:"id" bson:"id" validate:"required,numeric,gt=0"`
	Name string `json:"name" bson:"name" validate:"required,min=3"`
}

type User struct {
	ID 		 int `json:"id"`
	Name 	 string `json:"name" validate:"required,min=3"`
	Username string `json:"username" validate:"required,min=3"`
	Password string `json:"password" validate:"required,min=6"`
}

type CustomValidator struct {
	Validator *validator.Validate
}

func (c *CustomValidator) Validate(i interface{}) error {
	if err := c.Validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
