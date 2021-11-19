// Package models contains all models
package models

// Cats is simple cats
type Cats struct {
	ID   int32  `param:"id" json:"id" bson:"id" query:"id" header:"id" form:"id" xml:"id" validate:"required,numeric,gt=0"`
	Name string `param:"name" json:"name" bson:"name" query:"name" header:"name" form:"name" xml:"name" validate:"required,min=3"`
}

// User is simple user
type User struct {
	ID       int    `json:"id" param:"id" query:"id"`
	Name     string `json:"name" param:"name" query:"name" validate:"required,min=3"`
	Username string `json:"username" param:"username" query:"username" validate:"required,min=3"`
	Password string `json:"password" param:"password" query:"password" validate:"required,min=6"`
}
