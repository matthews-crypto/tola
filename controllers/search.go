package controllers

import (
    "context"
    "net/http"
    "tola/utils"
    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
)

func Search(c *gin.Context) {
    query := c.Query("query")
    if query == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "La requête de recherche ne peut pas être vide"})
        return
    }

    var users []struct {
        Name string `json:"name"`
    }
    var questions []struct {
        Title string `json:"title"`
    }
    var publications []struct {
        Title string `json:"title"`
    }

    userCursor, err := utils.UserCollection.Find(context.Background(), bson.M{"name": bson.M{"$regex": query, "$options": "i"}})
    if err == nil {
        userCursor.All(context.Background(), &users)
    }

    questionCursor, err := utils.QuestionCollection.Find(context.Background(), bson.M{"title": bson.M{"$regex": query, "$options": "i"}})
    if err == nil {
        questionCursor.All(context.Background(), &questions)
    }

    publicationCursor, err := utils.PublicationCollection.Find(context.Background(), bson.M{"title": bson.M{"$regex": query, "$options": "i"}})
    if err == nil {
        publicationCursor.All(context.Background(), &publications)
    }

    results := []struct {
        Title string `json:"title"`
        URL   string `json:"url"`
    }{}

    for _, user := range users {
        results = append(results, struct {
            Title string `json:"title"`
            URL   string `json:"url"`
        }{
            Title: user.Name,
            URL:   "/user/profile/" + user.Name,
        })
    }
    for _, question := range questions {
        results = append(results, struct {
            Title string `json:"title"`
            URL   string `json:"url"`
        }{
            Title: question.Title,
            URL:   "/questions/" + question.Title,
        })
    }
    for _, publication := range publications {
        results = append(results, struct {
            Title string `json:"title"`
            URL   string `json:"url"`
        }{
            Title: publication.Title,
            URL:   "/publications/" + publication.Title,
        })
    }

    c.JSON(http.StatusOK, gin.H{"results": results})
}
