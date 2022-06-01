package services

import (
	"testing"

	"github.com/croisade/chimichanga/models"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestAccountService(t *testing.T) {
	accountService := NewAccountServiceImpl(accountsCollection, ctx)
	t.Run("create Account", func(t *testing.T) {
		var got *models.Account
		account := &models.Account{AccountId: "123", Email: "test@example.com", Password: "password", FirstName: "first", LastName: "last"}

		accountService.CreateAccount(account)

		err := accountsCollection.FindOne(ctx, bson.M{"firstName": "first"}).Decode(&got)
		assert.Nil(t, err)
		assert.Equal(t, account, got)

	})

	t.Run("get Account", func(t *testing.T) {
		var got *models.Account
		want := &models.Account{AccountId: "123", Email: "test@example.com", Password: "password", FirstName: "first", LastName: "last"}

		got, err := accountService.GetAccount("123")

		assert.Nil(t, err)
		assert.Contains(t, got.AccountId, want.AccountId)
	})

	t.Run("get Accounts", func(t *testing.T) {
		var got []models.Account
		want := &models.Account{AccountId: "234", Email: "test@example.com", Password: "password", FirstName: "first", LastName: "last"}

		accountService.CreateAccount(want)
		got, err := accountService.GetAccounts()

		assert.Nil(t, err)
		assert.Equal(t, 2, len(got))

	})

	t.Run("Update Account", func(t *testing.T) {
		want := &models.Account{AccountId: "123", Email: "test@example.com", Password: "password", FirstName: "first", LastName: "last"}
		input := &models.Account{AccountId: "123", FirstName: "Middle"}

		accountService.UpdateAccount(input)
		got, err := accountService.GetAccount("123")

		assert.Nil(t, err)
		assert.NotEqual(t, want.FirstName, got.FirstName)
	})
}
