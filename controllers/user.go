package controllers

import (
	"context"
	"net/http"
	"tola/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func UpdateUserCategories(c *gin.Context) {
	userID := c.Param("id")
	var body struct {
		Categories []string `json:"categories"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(body.Categories) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Veuillez sélectionner au moins 3 catégories"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID utilisateur invalide"})
		return
	}

	_, err = utils.UserCollection.UpdateOne(context.Background(), bson.M{"_id": objID}, bson.M{
		"$set": bson.M{"categories": body.Categories},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour des catégories de l'utilisateur"})
		return
	}

	// Créer une session pour l'utilisateur
	session := sessions.Default(c)
	session.Set("userID", userID)
	session.Save()

	c.JSON(http.StatusOK, gin.H{"message": "Catégories mises à jour avec succès", "redirect": "/public/profile.html"})
}
func GetUserInfo(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID")
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non connecté"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID utilisateur invalide"})
		return
	}

	var user struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	err = utils.UserCollection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Utilisateur non trouvé"})
		return
	}

	c.JSON(http.StatusOK, user)
}
func GetUserPosts(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID")
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non connecté"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID utilisateur invalide"})
		return
	}

	var user struct {
		Categories []string `bson:"categories"`
	}
	err = utils.UserCollection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Utilisateur non trouvé"})
		return
	}

	var posts []struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	cursor, err := utils.QuestionCollection.Find(context.Background(), bson.M{
		"categories": bson.M{"$in": user.Categories},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des publications"})
		return
	}

	if err := cursor.All(context.Background(), &posts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la lecture des publications"})
		return
	}

	c.JSON(http.StatusOK, posts)
}
func UpdateUserProfile(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID")
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non connecté"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID utilisateur invalide"})
		return
	}

	var body struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := bson.M{}
	if body.Name != "" {
		update["name"] = body.Name
	}
	if body.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du hachage du mot de passe"})
			return
		}
		update["password"] = string(hashedPassword)
	}

	_, err = utils.UserCollection.UpdateOne(context.Background(), bson.M{"_id": objID}, bson.M{
		"$set": update,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour des informations de l'utilisateur"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Informations mises à jour avec succès"})
}
func GetUserProfile(c *gin.Context) {
	userID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID utilisateur invalide"})
		return
	}

	var user struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	err = utils.UserCollection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Utilisateur non trouvé"})
		return
	}

	var questions []struct {
		Title string `json:"title"`
		ID    string `json:"id"`
	}
	var publications []struct {
		Title string `json:"title"`
		ID    string `json:"id"`
	}

	questionCursor, err := utils.QuestionCollection.Find(context.Background(), bson.M{"user_id": objID})
	if err == nil {
		questionCursor.All(context.Background(), &questions)
	}

	publicationCursor, err := utils.PublicationCollection.Find(context.Background(), bson.M{"user_id": objID})
	if err == nil {
		publicationCursor.All(context.Background(), &publications)
	}

	c.JSON(http.StatusOK, gin.H{
		"user":         user,
		"questions":    questions,
		"publications": publications,
	})
}

func GetQuestionDetails(c *gin.Context) {
	questionID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(questionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de question invalide"})
		return
	}

	var question struct {
		Title    string `json:"title"`
		Content  string `json:"content"`
		UserID   string `json:"user_id"`
		Category string `json:"category"`
	}
	err = utils.QuestionCollection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&question)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Question non trouvée"})
		return
	}

	var answers []struct {
		Content string `json:"content"`
		UserID  string `json:"user_id"`
	}
	answerCursor, err := utils.AnswerCollection.Find(context.Background(), bson.M{"question_id": objID})
	if err == nil {
		answerCursor.All(context.Background(), &answers)
	}

	c.JSON(http.StatusOK, gin.H{
		"question": question,
		"answers":  answers,
	})
}
