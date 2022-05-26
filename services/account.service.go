package services

import (
	"context"

	"github.com/croisade/chimichanga/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AccountService interface {
	CreateAccount(*models.Account) (*models.Account, error)
	GetAccount(string) (*models.Account, error)
	GetAccounts() ([]*models.Account, error)
	DeleteAccount(string) error
	UpdateAccount(*models.Account) error
}

type AccountServiceImpl struct {
	accountcollection *mongo.Collection
	ctx               context.Context
}

func NewAccountServiceImpl(accountcollection *mongo.Collection, ctx context.Context) *AccountServiceImpl {
	return &AccountServiceImpl{
		accountcollection: accountcollection,
		ctx:               ctx,
	}
}

func (s *AccountServiceImpl) CreateAccount(account *models.Account) (*models.Account, error) {
	_, err := s.accountcollection.InsertOne(s.ctx, account)
	return account, err
}

func (s *AccountServiceImpl) GetAccount(accountId string) (*models.Account, error) {
	var result *models.Account
	var err error

	filter := bson.M{"accountId": accountId}
	err = s.accountcollection.FindOne(s.ctx, filter).Decode(&result)

	return result, err
}

func (s *AccountServiceImpl) GetAccounts(accountId string) ([]*models.Account, error) {
	var results []*models.Account
	var err error
	var cursor *mongo.Cursor
	cursor, err = s.accountcollection.Find(s.ctx, bson.M{})

	if err == nil {
		return results, err
	}

	err = cursor.All(s.ctx, results)

	return results, err
}

func (s *AccountServiceImpl) DeleteAccount(accountId string) error {
	_, err := s.accountcollection.DeleteOne(s.ctx, bson.M{"accountId": accountId})
	return err
}

func (s *AccountServiceImpl) UpdateAccount(account *models.Account) error {
	filter := bson.M{"accountId": account.AccountId}
	update := bson.M{"$set": bson.M{
		"AccountId": account.AccountId,
		"Email":     account.Email,
		"Password":  account.Password,
		"FirstName": account.FirstName,
		"LastName":  account.LastName,
	}}

	_, err := s.accountcollection.UpdateOne(s.ctx, filter, update)

	return err
}
