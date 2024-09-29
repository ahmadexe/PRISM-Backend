package repository

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/ahmadexe/prism-backend/services/jobs/data"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type JobsRepo struct {
	jobsCollection         *mongo.Collection
	applicationsCollection *mongo.Collection
}

func NewJobsRepo(client *mongo.Client) *JobsRepo {
	collection := client.Database("jobs-db").Collection("jobs")
	applicationsCollection := client.Database("jobs-db").Collection("applications")
	return &JobsRepo{jobsCollection: collection, applicationsCollection: applicationsCollection}
}

func (jr *JobsRepo) CreateJob(ctx *gin.Context, job *data.Job) {
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := jr.jobsCollection.InsertOne(context, job)
	if err != nil {
		log.Println(err)
		ctx.JSON(500, gin.H{"error": "Internal server error. Please try again later."})
		return
	}

	ctx.JSON(201, gin.H{"message": "Job created successfully."})
}

func (jr *JobsRepo) GetJobs(ctx *gin.Context) {
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := jr.jobsCollection.Find(context, bson.M{
		"hired": primitive.Null{},
	})
	if err != nil {
		log.Println(err)
		ctx.JSON(500, gin.H{"error": "Internal server error. Please try again later."})
		return
	}

	var jobs []data.Job
	if err = cursor.All(context, &jobs); err != nil {
		log.Println(err)
		ctx.JSON(500, gin.H{"error": "Internal server error. Please try again later."})
		return
	}

	ctx.JSON(200, gin.H{"data": jobs})
}

func (jr *JobsRepo) GetJob(ctx *gin.Context, id primitive.ObjectID) {
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var job data.Job

	err := jr.jobsCollection.FindOne(context, bson.M{"_id": id}).Decode(&job)
	if err != nil {
		log.Println(err)
		ctx.JSON(500, gin.H{"error": "Internal server error. Please try again later."})
		return
	}

	ctx.JSON(200, gin.H{"data": job})
}

func (jr *JobsRepo) UpdateJob(ctx *gin.Context, job *data.Job) {
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := jr.jobsCollection.UpdateOne(context, bson.M{"_id": job.ID}, bson.M{"$set": job})
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Internal server error. Please try again later."})
		return
	}

	ctx.JSON(200, gin.H{"message": "Job updated successfully."})
}

func (jr *JobsRepo) DeleteJob(ctx *gin.Context, id primitive.ObjectID) {
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := jr.jobsCollection.DeleteOne(context, bson.M{"_id": id})
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Internal server error. Please try again later."})
		return
	}

	ctx.JSON(200, gin.H{"message": "Job deleted successfully."})
}

func (jr *JobsRepo) ToggleLikeOnJob(ctx *gin.Context, id primitive.ObjectID, userId primitive.ObjectID) {
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var job data.Job
	filter := bson.M{"_id": id, "likedBy": userId}
	err := jr.jobsCollection.FindOne(context, filter).Decode(&job)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			_, err = jr.jobsCollection.UpdateOne(context, bson.M{"_id": id}, bson.M{"$push": bson.M{"likedBy": userId}})
			if err != nil {
				log.Println(err)
				ctx.JSON(500, gin.H{"error": "Internal server error. Please try again later."})
				return
			}
			ctx.JSON(200, gin.H{"message": "Job liked successfully."})
			return
		}
		log.Println(err)
		ctx.JSON(500, gin.H{"error": "Internal server error. Please try again later."})
		return
	}

	_, err = jr.jobsCollection.UpdateOne(context, bson.M{"_id": id}, bson.M{"$pull": bson.M{"likedBy": userId}})
	if err != nil {
		log.Println(err)
		ctx.JSON(500, gin.H{"error": "Internal server error. Please try again later."})
		return
	}

	ctx.JSON(200, gin.H{"message": "Job unliked successfully."})
}

func (jr *JobsRepo) ApplyForJob(ctx *gin.Context, application data.JobApplication) {
	context, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var applicationData data.JobApplication
	filter := bson.M{"jobId": application.JobId, "userId": application.UserId}
	err := jr.applicationsCollection.FindOne(context, filter).Decode(&applicationData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			_, err = jr.applicationsCollection.InsertOne(context, application)
			if err != nil {
				log.Println(err)
				ctx.JSON(500, gin.H{"error": "Internal server error. Please try again later."})
				return
			}

			_, err = jr.jobsCollection.UpdateOne(context, bson.M{"_id": application.JobId}, bson.M{"$push": bson.M{"appliedBy": application.UserId}})
			if err != nil {
				log.Println(err)
				ctx.JSON(500, gin.H{"error": "Internal server error. Please try again later."})
				return
			}

			ctx.JSON(201, gin.H{"message": "Application submitted successfully."})
			return
		}
		log.Println(err)
		ctx.JSON(500, gin.H{"error": "Internal server error. Please try again later."})
		return
	}

	ctx.JSON(400, gin.H{"error": "You have already applied for this job."})
}

func (jr *JobsRepo) HireForJob(ctx *gin.Context, id primitive.ObjectID, userId primitive.ObjectID) {
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	time := time.Now().UnixMicro()

	_, err := jr.jobsCollection.UpdateOne(context, bson.M{"_id": id}, bson.M{"$set": bson.M{"hired": userId, "hiredAt": time}})
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Internal server error. Please try again later."})
		return
	}

	_, err = jr.applicationsCollection.UpdateOne(context, bson.M{"jobId": id, "userId": userId}, bson.M{"$set": bson.M{"isHired": true}})
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Internal server error. Please try again later."})
		return
	}

	ctx.JSON(200, gin.H{"message": "User hired successfully."})
}

func (jr *JobsRepo) JobsByMe(ctx *gin.Context, id primitive.ObjectID) {
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := jr.jobsCollection.Find(context, bson.M{"postedBy": id})
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please try again later."})
		return
	}

	var jobs []data.Job

	if err = cursor.All(context, &jobs); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please try again later."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": jobs})
}

func (jr *JobsRepo) JobsLikedByMe(ctx *gin.Context, id primitive.ObjectID) {
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := jr.jobsCollection.Find(context, bson.M{"likedBy": id})
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please try again later."})
		return
	}

	var jobs []data.Job

	if err = cursor.All(context, &jobs); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please try again later."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": jobs})
}

func (jr *JobsRepo) JobsAppliedByMe(ctx *gin.Context, id primitive.ObjectID) {
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := jr.jobsCollection.Find(context, bson.M{"appliedBy": id})
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please try again later."})
		return
	}

	var jobs []data.Job

	if err = cursor.All(context, &jobs); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please try again later."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": jobs})
}
