package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Account struct {
	AccountId string              `json:"accountId" bson:"accountId"`
	Email     string              `json:"email" bson:"email" binding:"required"`
	Password  string              `json:"password" bson:"password" binding:"required"`
	FirstName string              `json:"firstName" bson:"firstName" binding:"required"`
	LastName  string              `json:"lastName" bson:"lastName" binding:"required"`
	CreatedAt primitive.Timestamp `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt primitive.Timestamp `json:"updatedAt" bson:"updatedAt,omitempty"`
}
