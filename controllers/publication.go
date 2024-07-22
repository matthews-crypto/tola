package controllers

import (
	"context"
	"net/http"
	"path/filepath"
	"time"
	"tola/models"
	"tola/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreatePublication(c *gin.Context) {
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

	var publication models.Publication
	publication.UserID = objID
	publication.CreatedAt = time.Now()
	publication.Likes = 0
	publication.Dislikes = 0

	// Récupérer les champs du formulaire
	publication.Title = c.PostForm("title")
	publication.Content = c.PostForm("content")
	publication.Category = c.PostForm("category")

	// Gestion du fichier image
	file, err := c.FormFile("image")
	if err == nil {
		filePath := filepath.Join("public/uploads", file.Filename)
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de l'enregistrement de l'image"})
			return
		}
		publication.Image = filePath
	}

	_, err = utils.PublicationCollection.InsertOne(context.Background(), publication)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de la publication"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Publication créée avec succès"})
}

func ListUserPublications(c *gin.Context) {
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

	var publications []bson.M
	cursor, err := utils.PublicationCollection.Find(
		context.Background(),
		bson.M{"user_id": objID},
		options.Find().SetSort(bson.D{{"created_at", -1}}),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des publications"})
		return
	}
	if err = cursor.All(context.Background(), &publications); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la lecture des publications"})
		return
	}

	for i, pub := range publications {
		if pub["created_at"] != nil {
			if t, ok := pub["created_at"].(primitive.DateTime); ok {
				publications[i]["created_at"] = t.Time()
			}
		}
		// Assurez-vous que _id est une chaîne de caractères
		if id, ok := pub["_id"].(primitive.ObjectID); ok {
			publications[i]["_id"] = id.Hex()
		}
	}

	c.JSON(http.StatusOK, publications)
}

func DeletePublication(c *gin.Context) {
	publicationID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(publicationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de publication invalide"})
		return
	}

	_, err = utils.PublicationCollection.DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression de la publication"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Publication supprimée avec succès"})
}

func LikePublication(c *gin.Context) {
	publicationID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(publicationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de publication invalide"})
		return
	}

	_, err = utils.PublicationCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": objID},
		bson.M{"$inc": bson.M{"likes": 1}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de l'ajout du like"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Publication likée avec succès"})
}

func DislikePublication(c *gin.Context) {
	publicationID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(publicationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de publication invalide"})
		return
	}

	_, err = utils.PublicationCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": objID},
		bson.M{"$inc": bson.M{"dislikes": 1}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de l'ajout du dislike"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Publication dislikée avec succès"})
}

func ListUserCategoryPublications(c *gin.Context) {
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

	var publications []bson.M
	cursor, err := utils.PublicationCollection.Find(
		context.Background(),
		bson.M{"category": bson.M{"$in": user.Categories}},
		options.Find().SetSort(bson.D{{"created_at", -1}}),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des publications"})
		return
	}
	if err = cursor.All(context.Background(), &publications); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la lecture des publications"})
		return
	}

	for i, pub := range publications {
		if pub["created_at"] != nil {
			if t, ok := pub["created_at"].(primitive.DateTime); ok {
				publications[i]["created_at"] = t.Time()
			}
		}
		// Assurez-vous que _id est une chaîne de caractères
		if id, ok := pub["_id"].(primitive.ObjectID); ok {
			publications[i]["_id"] = id.Hex()
		}
	}

	c.JSON(http.StatusOK, publications)
}
