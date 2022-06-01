package services

import (
	"context"
	"time"

	"github.com/croisade/chimichanga/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	account.CreatedAt = primitive.Timestamp{T: uint32(time.Now().Unix())}
	account.UpdatedAt = primitive.Timestamp{T: uint32(time.Now().Unix())}

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

func (s *AccountServiceImpl) GetAccounts() ([]*models.Account, error) {
	var results []*models.Account
	var err error
	var cursor *mongo.Cursor
	cursor, err = s.accountcollection.Find(s.ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	err = cursor.All(s.ctx, &results)
	return results, err
}

func (s *AccountServiceImpl) DeleteAccount(accountId string) error {
	_, err := s.accountcollection.DeleteOne(s.ctx, bson.M{"accountId": accountId})
	return err
}

func (s *AccountServiceImpl) UpdateAccount(account *models.Account) error {
	filter := bson.M{"accountId": account.AccountId}
	var err error

	existingAccount, err := s.GetAccount(account.AccountId)

	if err != nil {
		return err
	}

	if account.Email != "" {
		existingAccount.Email = account.Email
	}
	if account.Password != "" {
		existingAccount.Password = account.Password
	}
	if account.FirstName != "" {
		existingAccount.FirstName = account.FirstName
	}
	if account.LastName != "" {
		existingAccount.LastName = account.LastName
	}

	existingAccount.CreatedAt = account.CreatedAt
	existingAccount.UpdatedAt = primitive.Timestamp{T: uint32(time.Now().Unix())}

	_, err = s.accountcollection.UpdateOne(s.ctx, filter, bson.M{"$set": existingAccount})

	if err != nil {
		return err
	}

	return nil
}
