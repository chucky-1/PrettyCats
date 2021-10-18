package models


type Cats struct {
	ID   int32  `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}

//type CatsInter interface {
//	GetAllCats() ([]*Cats, error)
//	CreateCats(jsonMap map[string]interface{}) (*Cats, error)
//	GetCat(id string) (*Cats, error)
//	UpdateCat(id string, jsonMap map[string]interface{}) (*Cats, error)
//	DeleteCat(id string) (*Cats, error)
//}
