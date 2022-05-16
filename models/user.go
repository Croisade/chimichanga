package models

import "time"

type Runs struct {
	Pace      float64   `json:"pace" bson:"pace"`
	Time      string    `json:"time" bson:"time"`
	Distance  string    `json:"distance" bson:"distance"`
	Date      time.Time `json:"date" bson:"date"`
	Lap       int       `json:"lap" bson:"lap"`
	SessionId string    `json:"sessionId" bson:"sessionId"`
	UserId    string    `json:"userId" bson:"userId"`
}

type User struct {
	UserId string `json:"userId" bson:"userId"`
	Runs   Runs   `json:"runs" bson:"runs"`
}
