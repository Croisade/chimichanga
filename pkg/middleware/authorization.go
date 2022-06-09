package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/croisade/chimichanga/pkg/services"
	"github.com/gin-gonic/gin"
)

func AuthorizeUserJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		const BEARER_SCHEMA = "Bearer "
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": "request missing token"})
			return
		}

		tokenString := authHeader[len(BEARER_SCHEMA):]
		token, err := services.NewJWTAuthService().ValidateToken(tokenString)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": err.Error()})
			return
		}

		// https://github.com/dgrijalva/jwt-go/issues/262
		claims := services.MyCustomClaims{}
		tmp, _ := json.Marshal(token.Claims)
		_ = json.Unmarshal(tmp, &claims)

		if claims.Group != "USER" && claims.Group != "ADMIN" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": "invalid group"})
			return
		}

		ctx.Next()
	}
}

func AuthorizeAdminJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		const BEARER_SCHEMA = "Bearer "
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": "request missing token"})
			return
		}
		tokenString := authHeader[len(BEARER_SCHEMA):]
		token, err := services.NewJWTAuthService().ValidateToken(tokenString)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": err.Error()})
			return
		}

		// https://github.com/dgrijalva/jwt-go/issues/262
		claims := services.MyCustomClaims{}
		tmp, _ := json.Marshal(token.Claims)
		_ = json.Unmarshal(tmp, &claims)

		if claims.Group != "ADMIN" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": "invalid group"})
			return
		}

		ctx.Next()
	}
}
