package controllers

import (
    "context"
    "net/http"
    "tola/utils"
    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
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

    c.JSON(http.StatusOK, gin.H{"message": "Catégories mises à jour avec succès"})
}
