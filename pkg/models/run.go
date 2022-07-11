package models

import (
	"time"
)

type Run struct {
	Speed     float32   `json:"speed,omitempty" bson:"speed,omitempty"`
	Time      string    `json:"time,omitempty" bson:"time,omitempty"`
	Distance  float32   `json:"distance,omitempty" bson:"distance,omitempty"`
	Lap       int       `json:"lap,omitempty" bson:"lap,omitempty"`
	Incline   float32   `json:"incline,omitempty" bson:"incline,omitempty"`
	RunId     string    `json:"runId,omitempty" bson:"runId,omitempty"`
	AccountId string    `json:"accountId,omitempty" bson:"accountId,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}
