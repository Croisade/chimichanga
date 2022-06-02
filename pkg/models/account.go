package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Account struct {
	AccountId string              `json:"accountId" bson:"accountId"`
	Email     string              `json:"email" bson:"email"`
	Password  string              `json:"password" bson:"password"`
	FirstName string              `json:"firstName" bson:"firstName"`
	LastName  string              `json:"lastName" bson:"lastName"`
	CreatedAt primitive.Timestamp `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt primitive.Timestamp `json:"updatedAt" bson:"updatedAt,omitempty"`
}
