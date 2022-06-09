package services

import (
	"testing"

	"github.com/croisade/chimichanga/pkg/models"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestAccountService(t *testing.T) {
	accountService := NewAccountServiceImpl(accountsCollection, ctx)
	want := &models.Account{Email: "test@example.com", Password: "password", FirstName: "first", LastName: "last"}

	t.Run("create Account", func(t *testing.T) {
		var got *models.Account
		accountService.CreateAccount(want)

		err := accountsCollection.FindOne(ctx, bson.M{"firstName": "first"}).Decode(&got)

		want.AccountId = got.AccountId

		assert.Nil(t, err)
		assert.Equal(t, want, got)

	})

	t.Run("get Account", func(t *testing.T) {
		var got *models.Account
		got, err := accountService.GetAccount(want.AccountId)

		assert.Nil(t, err)
		assert.Contains(t, got.AccountId, want.AccountId)
	})

	t.Run("get Accounts", func(t *testing.T) {
		var got []*models.Account
		want := &models.Account{Email: "test@example.com", Password: "password", FirstName: "first", LastName: "last"}

		accountService.CreateAccount(want)
		got, err := accountService.GetAccounts()

		assert.Nil(t, err)
		assert.Equal(t, 2, len(got))

	})

	t.Run("Update Account", func(t *testing.T) {
		input := &models.Account{AccountId: want.AccountId, FirstName: "Middle"}

		got, err := accountService.UpdateAccount(input)

		assert.Nil(t, err)
		assert.NotEqual(t, want.FirstName, got.FirstName)
	})

	t.Run("Delete Account", func(t *testing.T) {
		accountService.DeleteAccount(want.AccountId)
		_, err := accountService.GetAccount("123")

		assert.ErrorContains(t, err, "no documents")
	})
}
