// Package models contains all models
package models

// Cat is simple cats
type Cat struct {
	ID   int32  `param:"id" json:"id" bson:"id" query:"name" validate:"required,numeric,gt=0"`
	Name string `json:"name" bson:"name" query:"name" validate:"required,min=3"`
}

// User is simple user
type User struct {
	ID       int    `param:"id" json:"id" query:"id"`
	Name     string `json:"name" query:"name" validate:"required,min=3"`
	Username string `json:"username" query:"username" validate:"required,min=3"`
	Password string `json:"password" query:"password" validate:"required,min=6"`
}
