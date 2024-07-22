package controllers

import (
	"context"
	"net/http"
	"time"
	"tola/models"
	"tola/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AskQuestion(c *gin.Context) {
	var question models.Question
	if err := c.ShouldBindJSON(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	question.ID = primitive.NewObjectID()
	question.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	_, err := utils.QuestionCollection.InsertOne(context.Background(), question)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de la question"})
		return
	}

	c.JSON(http.StatusOK, question)
}

func ListQuestions(c *gin.Context) {
	var questions []models.Question
	cursor, err := utils.QuestionCollection.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des questions"})
		return
	}

	if err := cursor.All(context.Background(), &questions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la lecture des questions"})
		return
	}

	c.JSON(http.StatusOK, questions)
}

func AnswerQuestion(c *gin.Context) {
	var answer models.Answer
	if err := c.ShouldBindJSON(&answer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	answer.ID = primitive.NewObjectID()
	answer.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	_, err := utils.AnswerCollection.InsertOne(context.Background(), answer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de la réponse"})
		return
	}

	c.JSON(http.StatusOK, answer)
}

// func ListAnswers(c *gin.Context) {
// 	questionID := c.Param("questionId")
// 	objID, err := primitive.ObjectIDFromHex(questionID)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de question invalide"})
// 		return
// 	}

// 	var answers []models.Answer
// 	cursor, err := utils.AnswerCollection.Find(context.Background(), bson.M{"questionId": objID})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des réponses"})
// 		return
// 	}

// 	if err := cursor.All(context.Background(), &answers); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la lecture des réponses"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, answers)
// }
