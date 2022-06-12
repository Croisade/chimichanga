package services

import (
	"testing"

	"github.com/croisade/chimichanga/pkg/models"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func createAccountRoutine(f *AccountServiceImpl, list []*models.Account, out chan<- []*models.Account) {
	var result []*models.Account
	for i := 0; i < 2; i++ {
		account, _ := f.CreateAccount(list[i])
		result = append(result, account)
	}
	out <- result
}

func TestAccountService(t *testing.T) {
	accountService := NewAccountServiceImpl(accountCollection, ctx)
	want := &models.Account{Email: "test@example.com", Password: "password", FirstName: "first", LastName: "last"}

	t.Run("create Account", func(t *testing.T) {
		accountCollection.DeleteMany(ctx, bson.D{{}})
		var got *models.Account
		accountService.CreateAccount(want)

		err := accountCollection.FindOne(ctx, bson.M{"firstName": "first"}).Decode(&got)

		want.AccountId = got.AccountId

		assert.Nil(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("Should error out if account already exists", func(t *testing.T) {
		accountCollection.DeleteMany(ctx, bson.D{{}})
		accountService.CreateAccount(want)
		_, err := accountService.CreateAccount(want)

		assert.NotNil(t, err)
	})

	t.Run("get Account", func(t *testing.T) {
		accountCollection.DeleteMany(ctx, bson.D{{}})
		_, err := accountService.CreateAccount(want)
		var got *models.Account
		got, err = accountService.GetAccount(want.AccountId)

		assert.Nil(t, err)
		assert.Contains(t, got.AccountId, want.AccountId)
	})

	t.Run("Update Account", func(t *testing.T) {
		accountCollection.DeleteMany(ctx, bson.D{{}})
		_, err := accountService.CreateAccount(want)
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

	t.Run("get Accounts", func(t *testing.T) {
		outSlice := make(chan []*models.Account)
		out := make(chan *mongo.DeleteResult)
		go func() {
			result, _ := accountCollection.DeleteMany(ctx, bson.D{{}})
			out <- result
		}()
		<-out
		var got []*models.Account

		firstAcc := &models.Account{Email: "test3@example.com", Password: "password", FirstName: "first", LastName: "last"}
		newAcc := &models.Account{Email: "test2@example.com", Password: "password", FirstName: "first", LastName: "last"}

		go createAccountRoutine(accountService, []*models.Account{firstAcc, newAcc}, outSlice)
		<-outSlice

		got, err := accountService.GetAccounts()
		// man, _ := json.MarshalIndent(got, "", "    ")
		// fmt.Println(string(man))

		assert.Nil(t, err)
		assert.Equal(t, 2, len(got))
		// assert.Equal(t, got[0].AccountId, .AccountId)
	})
}
