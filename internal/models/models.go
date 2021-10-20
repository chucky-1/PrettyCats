package models


type Cats struct {
	ID   int32  `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}
