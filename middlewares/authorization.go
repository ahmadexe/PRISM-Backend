package middlewares

import (
	"net/http"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
)

type authMiddleware struct {
	App *firebase.App
}

func InitAuthMiddleware(app *firebase.App) *authMiddleware {
	return &authMiddleware{App: app}
}

func (middleware *authMiddleware) VerifyUser(ctx *gin.Context) {
	authToken := ctx.GetHeader("Authorization")

	client, err := middleware.App.Auth(ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please try again later."})

		return
	}

	_, err = client.VerifyIDToken(ctx, authToken)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token. Please login again."})

		return
	}
}