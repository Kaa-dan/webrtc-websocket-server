package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID         primitive.ObjectID `bson:"_id"`
	First_name *string            `json:"first_name"`
	User_id    string             `json:"user_id"`
}
