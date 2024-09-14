package data

import (
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Job struct {
	ID          primitive.ObjectID   `json:"id" bson:"_id"`
	PostedBy    primitive.ObjectID   `json:"postedBy" bson:"postedBy" validate:"required"`
	Title       string               `json:"title" bson:"title" validate:"required"`
	Description string               `json:"description" bson:"description" validate:"required"`
	PostedAt    int32                `json:"postedAt" bson:"postedAt" validate:"required"`
	Country     string               `json:"country" bson:"country" validate:"required"`
	Keywords    []string             `json:"keywords" bson:"keywords" validate:"required"`
	LikedBy     []primitive.ObjectID `json:"likedBy" bson:"likedBy"`
	AppliedBy   []primitive.ObjectID `json:"appliedBy" bson:"appliedBy"`
	Hired       *primitive.ObjectID  `json:"hired" bson:"hired"`
	HiredAt     int64                `json:"hiredAt" bson:"hiredAt"`
	Budget      float64              `json:"budget" bson:"budget" validate:"required"`
	BudgetMeta  string               `json:"budgetMeta" bson:"budgetMeta" validate:"required"`
	Username    string               `json:"username" bson:"username" validate:"required"`
	Avatar      string               `json:"avatar" bson:"avatar"`
}

func (job *Job) Validate() error {
	v := validator.New()
	return v.Struct(job)
}
