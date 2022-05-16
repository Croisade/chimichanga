package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Runs struct {
	Pace      float64             `json:"pace" bson:"pace"`
	Time      string              `json:"time" bson:"time"`
	Distance  string              `json:"distance" bson:"distance"`
	Date      primitive.Timestamp `json:"date" bson:"date"`
	Lap       int                 `json:"lap" bson:"lap"`
	SessionId string              `json:"sessionId" bson:"sessionId"`
	UserId    string              `json:"userId" bson:"userId"`
}

type User struct {
	UserId    string              `json:"userId" bson:"userId"`
	CreatedAt primitive.Timestamp `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt primitive.Timestamp `json:"updatedAt" bson:"updatedAt,omitempty"`
	Runs      Runs                `json:"runs" bson:"runs"`
}
