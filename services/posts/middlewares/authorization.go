package middlewares

import (
	"context"
	"net/http"
	"strings"
	"time"

	pb "github.com/ahmadexe/prism-backend/grpc/auth/generated"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func VerifyUser(ctx *gin.Context) {
	authToken := ctx.GetHeader("Authorization")

	con, err := grpc.Dial(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please try again later."})

		return
	}
	defer con.Close()

	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	authToken = strings.TrimPrefix(authToken, "Bearer ")
	client := pb.NewAuthClient(con)

	res, err := client.Authorize(c, &pb.AuthorizeRequest{
		Token: authToken,
	})


	if err != nil {
		if strings.Contains(err.Error(), "invalid toke") {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "User is not authorized. Please login again."})

			return
		}

		ctx.AbortWithStatus(http.StatusInternalServerError)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please try again later."})
		return
	}

	if res.IsAuthorized {
		ctx.Next()
		return
	}
}
