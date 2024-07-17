package controllers

import (
    "context"
    "net/http"
    "tola/models"
    "tola/utils"
    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

func AddCategory(c *gin.Context) {
    var category models.Category
    if err := c.ShouldBindJSON(&category); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    category.ID = primitive.NewObjectID()
    _, err := utils.CategoryCollection.InsertOne(context.Background(), category)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de l'insertion de la catégorie"})
        return
    }

    c.JSON(http.StatusOK, category)
}

func ListCategories(c *gin.Context) {
    var categories []models.Category
    cursor, err := utils.CategoryCollection.Find(context.Background(), bson.M{})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des catégories"})
        return
    }

    if err := cursor.All(context.Background(), &categories); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la lecture des catégories"})
        return
    }

    c.JSON(http.StatusOK, categories)
}
