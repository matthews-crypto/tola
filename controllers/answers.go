package controllers

import (
	"context"
	"net/http"
	"time"
	"tola/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Answer struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	QuestionID primitive.ObjectID `bson:"question_id,omitempty"`
	UserID     primitive.ObjectID `bson:"user_id,omitempty"`
	Content    string             `bson:"content"`
	CreatedAt  time.Time          `bson:"created_at"`
	Likes      int                `bson:"likes"`
}

func CreateAnswer(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID")
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non connecté"})
		return
	}

	objUserID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID utilisateur invalide"})
		return
	}

	questionID := c.Param("questionId")
	objQuestionID, err := primitive.ObjectIDFromHex(questionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID question invalide"})
		return
	}

	var answer Answer
	if err := c.BindJSON(&answer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	answer.ID = primitive.NewObjectID()
	answer.QuestionID = objQuestionID
	answer.UserID = objUserID
	answer.CreatedAt = time.Now()
	answer.Likes = 0

	_, err = utils.AnswerCollection.InsertOne(context.Background(), answer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de la réponse"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Réponse créée avec succès"})
}

func ListAnswers(c *gin.Context) {
	questionID := c.Param("questionId")
	objQuestionID, err := primitive.ObjectIDFromHex(questionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID question invalide"})
		return
	}

	var answers []Answer
	cursor, err := utils.AnswerCollection.Find(context.Background(), bson.M{"question_id": objQuestionID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des réponses"})
		return
	}

	if err = cursor.All(context.Background(), &answers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la lecture des réponses"})
		return
	}

	c.JSON(http.StatusOK, answers)
}

func LikeAnswer(c *gin.Context) {
	answerID := c.Param("answerId")
	objAnswerID, err := primitive.ObjectIDFromHex(answerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID réponse invalide"})
		return
	}

	_, err = utils.AnswerCollection.UpdateOne(context.Background(), bson.M{"_id": objAnswerID}, bson.M{"$inc": bson.M{"likes": 1}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du like de la réponse"})
		return
	}

	var answer Answer
	err = utils.AnswerCollection.FindOne(context.Background(), bson.M{"_id": objAnswerID}).Decode(&answer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération de la réponse"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"likes": answer.Likes})
}
