package repositories

import (
	"context"
	"log"
	"net/http"
	"time"

	"golang.org/x/exp/slices"

	"github.com/ahmadexe/prism-backend/services/auth/data"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthRepo struct {
	collection *mongo.Collection
}

func InitAuthRepo(client *mongo.Client) *AuthRepo {
	collection := client.Database("auth-db").Collection("users")
	return &AuthRepo{collection: collection}
}

func (repo *AuthRepo) AddUser(user data.AuthData, ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := repo.collection.InsertOne(c, user)

	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding user to database. Please try again later."})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

func (repo *AuthRepo) GetUserById(ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	id := ctx.Param("id")
	filter := bson.M{"uid": id}
	var user data.AuthData

	if err := repo.collection.FindOne(c, filter).Decode(&user); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please try again later."})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (repo *AuthRepo) GetById(id primitive.ObjectID, ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"_id": id}
	var user data.AuthData

	if err := repo.collection.FindOne(c, filter).Decode(&user); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please try again later."})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (repo *AuthRepo) UpdateUser(user data.AuthData, ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": user.Id}
	update := bson.M{"$set": user}

	if _, err := repo.collection.UpdateOne(c, filter, update); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding user to database. Please try again later."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (repo *AuthRepo) ToggleFollow(req data.FollowRequest, ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": req.From}
	var fromUser data.AuthData

	if err := repo.collection.FindOne(c, filter).Decode(&fromUser); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please try again later."})
		return
	}

	if slices.Contains(fromUser.Following, req.To) {
		update := bson.M{"$pull": bson.M{"following": req.To}}
		update2 := bson.M{"$pull": bson.M{"followers": req.From}}

		filter2 := bson.M{"_id": req.To}

		if _, err := repo.collection.UpdateOne(c, filter, update); err != nil {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user in database. Please try again later."})
			return
		}

		if _, err := repo.collection.UpdateOne(c, filter2, update2); err != nil {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user in database. Please try again later."})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "User unfollowed successfully."})
	} else {
		update := bson.M{"$push": bson.M{"following": req.To}}
		update2 := bson.M{"$push": bson.M{"followers": req.From}}

		filter2 := bson.M{"_id": req.To}

		if _, err := repo.collection.UpdateOne(c, filter, update); err != nil {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user in database. Please try again later."})
			return
		}

		if _, err := repo.collection.UpdateOne(c, filter2, update2); err != nil {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user in database. Please try again later."})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "User followed successfully."})
	}
}

func (repo *AuthRepo) GetUserBySubString(sub string) ([]data.AuthData, error) {
	c, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"fullname": primitive.Regex{Pattern: sub, Options: "i"}}
	var users []data.AuthData

	cursor, err := repo.collection.Find(c, filter)
	if err != nil {
		return []data.AuthData{}, err
	}

	if err := cursor.All(c, &users); err != nil {
		return []data.AuthData{}, err
	}

	return users, nil
}

func (repo *AuthRepo) ToggleIsServiceProvider (id primitive.ObjectID, ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	
	var user data.AuthData

	if err := repo.collection.FindOne(c, filter).Decode(&user); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please try again later."})
		return
	}

	update := bson.M{"$set": bson.M{"isServiceProvider": !user.IsServiceProvider}}

	if _, err := repo.collection.UpdateOne(c, filter, update); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user in database. Please try again later."})
		return
	}

	user.IsServiceProvider = !user.IsServiceProvider

	ctx.JSON(http.StatusOK, gin.H{"message": "User is updated.", "data": user})
}

func (repo *AuthRepo) UpdateDeviceToken(ctx *gin.Context, tokenReq data.TokenRequest) {
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": tokenReq.UserId}

	update := bson.M{"$set": bson.M{"deviceToken": tokenReq.DeviceToken}}

	if _, err := repo.collection.UpdateOne(context, filter, update); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user in database. Please try again later."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Device token updated successfully."})
}