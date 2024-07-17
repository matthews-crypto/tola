package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Question struct {
    ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    UserID     primitive.ObjectID `bson:"userId" json:"userId"`
    Title      string             `bson:"title" json:"title"`
    Content    string             `bson:"content" json:"content"`
    Categories []string           `bson:"categories" json:"categories"`
    CreatedAt  primitive.DateTime `bson:"createdAt" json:"createdAt"`
}

type Answer struct {
    ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    QuestionID primitive.ObjectID `bson:"questionId" json:"questionId"`
    UserID     primitive.ObjectID `bson:"userId" json:"userId"`
    Content    string             `bson:"content" json:"content"`
    CreatedAt  primitive.DateTime `bson:"createdAt" json:"createdAt"`
}
