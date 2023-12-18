package models

type User struct {
	Age     int      `json:"age" bson:"age"`
	Name    string   `json:"name" bson:"name"`
	Id      string   `json:"id" bson:"_id"`
	Friends []string `json:"friends" bson:"friends, omitempty"`
}
