package models

import (
	"github.com/golang-jwt/jwt"
)


type Cats struct {
	ID   int32  `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}

type User struct {
	ID 		 int `json:"id"`
	Name 	 string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type JwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}