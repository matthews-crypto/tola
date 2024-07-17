package utils

import "go.mongodb.org/mongo-driver/mongo"

var UserCollection *mongo.Collection
var CategoryCollection *mongo.Collection
var QuestionCollection *mongo.Collection
var AnswerCollection *mongo.Collection

func InitCollections() {
	UserCollection = GetCollection("users")
	CategoryCollection = GetCollection("categories")
	QuestionCollection = GetCollection("questions")
	AnswerCollection = GetCollection("answers")
}
