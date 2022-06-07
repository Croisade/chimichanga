package controllers

import (
	"errors"
	"net/http"

	"github.com/croisade/chimichanga/pkg/models"
	"github.com/croisade/chimichanga/pkg/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AccountController struct {
	AccountService services.AccountService
	JWTService     services.JWTAuthService
}

type JWTtoken struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

type ErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func NewAccountController(accountService services.AccountService, jwtService services.JWTAuthService) AccountController {
	return AccountController{
		AccountService: accountService,
		JWTService:     jwtService,
	}
}

func (ac *AccountController) getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "lte":
		return "Should be less than " + fe.Param()
	case "gte":
		return "Should be greater than " + fe.Param()
	}
	return "Unknown error"
}

func (ac *AccountController) handleValidationError(ctx *gin.Context, err error) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {

		out := make([]ErrorMsg, len(ve))
		for i, fe := range ve {
			out[i] = ErrorMsg{fe.Field(), ac.getErrorMsg(fe)}
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": out})
		return
	}
}

func (ac *AccountController) CreateAccount(ctx *gin.Context) {
	var account models.Account
	if err := ctx.ShouldBindJSON(&account); err != nil {
		ac.handleValidationError(ctx, err)
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
	var accounts []*models.Account

	accounts, err := ac.AccountService.GetAccounts()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, accounts)
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

func (ac *AccountController) Login(ctx *gin.Context) {

	var login *services.Login

	if err := ctx.ShouldBindJSON(&login); err != nil {
		ac.handleValidationError(ctx, err)
	}

	_, err := ac.AccountService.Login(login)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
	}

	token, err := ac.JWTService.CreateToken()
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
	}

	refreshTokens, err := ac.JWTService.CreateRefreshToken()
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
	}

	result := JWTtoken{Token: token, RefreshToken: refreshTokens}
	ctx.JSON(http.StatusOK, result)
	return
}

func (ac *AccountController) Token(ctx *gin.Context)  {}
func (ac *AccountController) Logout(ctx *gin.Context) {}

func (ac *AccountController) RegisterAccountRoutes(rg *gin.RouterGroup) {
	accountRoute := rg.Group("/account")
	accountRoute.POST("/create", ac.CreateAccount)
	accountRoute.GET("/get/:accountId", ac.GetAccount)
	accountRoute.GET("/fetch", ac.GetAccounts)
	accountRoute.DELETE("/delete/:accountId", ac.DeleteAccount)
	accountRoute.PUT("/update", ac.UpdateAccount)
	accountRoute.PUT("/login", ac.UpdateAccount)
	accountRoute.PUT("/token", ac.UpdateAccount)
	accountRoute.PUT("/logout", ac.UpdateAccount)
}
