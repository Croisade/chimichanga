package controllers

import (
	"errors"
	"net/http"

	"github.com/croisade/chimichanga/pkg/middleware"
	"github.com/croisade/chimichanga/pkg/models"
	"github.com/croisade/chimichanga/pkg/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AccountController struct {
	AccountService services.AccountService
	JWTService     services.JWTAuthService
}

type RefreshToken struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type JWTtoken struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

type ErrorValidationMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Errors string `json:"errors"`
}

type UpdateAccountRequest struct {
	AccountId string `json:"accountId" bson:"accountId"`
	Email     string `json:"email" bson:"email"`
	Password  string `json:"password" bson:"password"`
	FirstName string `json:"firstName" bson:"firstName"`
	LastName  string `json:"lastName" bson:"lastName"`
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

		out := make([]ErrorValidationMsg, len(ve))
		for i, fe := range ve {
			out[i] = ErrorValidationMsg{fe.Field(), ac.getErrorMsg(fe)}
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": out})
		return
	}
	ctx.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
	return
}

func (ac *AccountController) CreateAccount(ctx *gin.Context) {
	var account models.Account
	if err := ctx.ShouldBindJSON(&account); err != nil {
		ac.handleValidationError(ctx, err)
		return
	}
	result, err := ac.AccountService.CreateAccount(&account)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"errors": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
	return
}

func (ac *AccountController) GetAccount(ctx *gin.Context) {
	accountId := ctx.Param("accountId")
	result, err := ac.AccountService.GetAccount(accountId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"errors": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
	return
}

func (ac *AccountController) GetAccounts(ctx *gin.Context) {
	var accounts []*models.Account

	accounts, err := ac.AccountService.GetAccounts()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"errors": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, accounts)
	return
}

func (ac *AccountController) DeleteAccount(ctx *gin.Context) {
	accountId := ctx.Param("accountId")
	err := ac.AccountService.DeleteAccount(accountId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"errors": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"errors": "success"})
	return
}

func (ac *AccountController) UpdateAccount(ctx *gin.Context) {
	var account *UpdateAccountRequest
	var accountToBeUpdated models.Account
	if err := ctx.ShouldBindJSON(&account); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}
	accountToBeUpdated.AccountId = account.AccountId
	accountToBeUpdated.Email = account.Email
	accountToBeUpdated.FirstName = account.FirstName
	accountToBeUpdated.LastName = account.LastName
	accountToBeUpdated.Password = account.Password

	result, err := ac.AccountService.UpdateAccount(&accountToBeUpdated)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"errors": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
	return
}

func (ac *AccountController) Login(ctx *gin.Context) {
	var login *services.LoginValidation

	if err := ctx.ShouldBindJSON(&login); err != nil {
		ac.handleValidationError(ctx, err)
		return
	}

	account, err := ac.AccountService.Login(login)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"errors": err.Error()})
		return
	}

	token, err := ac.JWTService.CreateToken()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
		return
	}

	refreshToken, err := ac.JWTService.CreateRefreshToken()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
		return
	}

	account.RefreshToken = refreshToken
	ac.AccountService.UpdateAccount(account)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"errors": err.Error()})
		return
	}

	result := JWTtoken{Token: token, RefreshToken: refreshToken}
	ctx.JSON(http.StatusOK, result)
	return
}

func (ac *AccountController) Token(ctx *gin.Context) {
	var refreshToken *RefreshToken
	if err := ctx.ShouldBindJSON(&refreshToken); err != nil {
		//@ @TODO Validate for jwt not for string
		ac.handleValidationError(ctx, err)
		return
	}

	_, err := ac.JWTService.ValidateToken(refreshToken.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
		return
	}

	token, err := ac.JWTService.CreateToken()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
		return
	}

	refreshTokens, err := ac.JWTService.CreateRefreshToken()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
		return
	}

	account, err := ac.AccountService.FindByRefreshToken(refreshToken.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
		return
	}

	_, err = ac.AccountService.UpdateAccount(account)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
		return
	}

	result := JWTtoken{Token: token, RefreshToken: refreshTokens}

	ctx.JSON(http.StatusOK, result)
	return
}
func (ac *AccountController) Logout(ctx *gin.Context) {
	var logoutValidation *services.LogoutValidation
	if err := ctx.ShouldBindJSON(&logoutValidation); err != nil {
		ac.handleValidationError(ctx, err)
		return
	}

	err := ac.AccountService.Logout(logoutValidation)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func (ac *AccountController) RegisterAccountRoutes(rg *gin.RouterGroup) {
	accountRouteNoMw := rg.Group("/account")
	accountRouteNoMw.POST("/create", ac.CreateAccount)
	accountRouteNoMw.PUT("/login", ac.Login)
	accountRouteUser := rg.Group("/account", middleware.AuthorizeUserJWT())
	accountRouteUser.GET("/get/:accountId", ac.GetAccount)
	accountRouteUser.DELETE("/delete/:accountId", ac.DeleteAccount)
	accountRouteUser.PUT("/update", ac.UpdateAccount)
	accountRouteUser.PUT("/token", ac.Token)
	accountRouteUser.PUT("/logout", ac.Logout)
	accountRouteAdmin := rg.Group("/account", middleware.AuthorizeAdminJWT())
	accountRouteAdmin.GET("/fetch", ac.GetAccounts)
}
