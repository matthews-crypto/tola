package controllers

import (
	"context"
	"net/http"
	"strings"
	"tola/models"
	"tola/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Vérifier si l'email se termine par @esp.sn
	if !strings.HasSuffix(user.Email, "@esp.sn") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "L'inscription est réservée aux étudiants de l'ESP"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du hachage du mot de passe"})
		return
	}
	user.Password = string(hashedPassword)
	user.Categories = []string{} // Initialiser la liste des catégories comme vide

	user.ID = primitive.NewObjectID()
	_, err = utils.UserCollection.InsertOne(context.Background(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de l'insertion de l'utilisateur"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Inscription réussie", "userId": user.ID.Hex()})
}

func Login(c *gin.Context) {
    var input models.User
    var user models.User
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    err := utils.UserCollection.FindOne(context.Background(), bson.M{"email": input.Email}).Decode(&user)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou mot de passe incorrect"})
        return
    }

    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou mot de passe incorrect"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Connexion réussie", "redirect": "/public/profile.html"})
}