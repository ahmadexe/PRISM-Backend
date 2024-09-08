package repository

import "go.mongodb.org/mongo-driver/mongo"

type JobsRepo struct {
	jobsCollection *mongo.Collection
}

func NewJobsRepo(client *mongo.Client) *JobsRepo {
	collection := client.Database("jobs-db").Collection("jobs")
	return &JobsRepo{jobsCollection: collection}
}