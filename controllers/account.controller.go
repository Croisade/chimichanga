package controllers

import (
	"net/http"

	"github.com/croisade/chimichanga/models"
	"github.com/croisade/chimichanga/services"
	"github.com/gin-gonic/gin"
)

type AccountController struct {
	AccountService services.AccountService
}

func NewAccountController(accountService services.AccountService) AccountController {
	return AccountController{
		AccountService: accountService,
	}
}

func (ac *AccountController) CreateAccount(ctx *gin.Context) {
	var account models.Account
	if err := ctx.ShouldBindJSON(&account); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	result, err := ac.AccountService.CreateAccount(&account)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
	return
}

func (ac *AccountController) GetAccount(ctx *gin.Context) {
	accountId := ctx.Param("AccountId")
	result, err := ac.AccountService.GetAccount(accountId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
	return
}

func (ac *AccountController) GetAccounts(ctx *gin.Context) {
	var users []*models.Account

	users, err := ac.AccountService.GetAccounts()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, users)
	return
}

func (ac *AccountController) DeleteAccount(ctx *gin.Context) {
	accountId := ctx.Param("AccountId")
	err := ac.AccountService.DeleteAccount(accountId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	return
}

func (ac *AccountController) UpdateAccount(ctx *gin.Context) {
	var account models.Account
	if err := ctx.ShouldBindJSON(&account); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	result, err := ac.AccountService.CreateAccount(&account)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
	return
}

func (ac *AccountController) RegisterAccountRoutes(rg *gin.RouterGroup) {
	accountRoute := rg.Group("/account")
	accountRoute.POST("/create", ac.CreateAccount)
	accountRoute.GET("/get/:name", ac.GetAccount)
	accountRoute.GET("/fetch", ac.GetAccounts)
	accountRoute.DELETE("/delete/:userId", ac.DeleteAccount)
	accountRoute.PUT("/update", ac.UpdateAccount)
}
