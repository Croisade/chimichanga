package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Run struct {
	Pace      float32             `json:"pace" bson:"pace"`
	Time      string              `json:"time" bson:"time"`
	Distance  float32             `json:"distance" bson:"distance"`
	Lap       int                 `json:"lap" bson:"lap"`
	Incline   float32             `json:"incline" bson:"incline"`
	RunId     string              `json:"runId" bson:"runId"`
	AccountId string              `json:"accountId" bson:"accountId"`
	CreatedAt primitive.Timestamp `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt primitive.Timestamp `json:"updatedAt" bson:"updatedAt,omitempty"`
}
