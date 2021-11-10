package models

type Cats struct {
	ID   int32  `json:"id" bson:"id" validate:"required,numeric,gt=0"`
	Name string `json:"name" bson:"name" validate:"required,min=3"`
}

type User struct {
	ID 		 int `json:"id" validate:"required"`
	Name 	 string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,gt=5"`
}
