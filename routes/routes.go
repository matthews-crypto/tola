package routes

import (
    "github.com/gin-gonic/gin"
    "tola/controllers"
    "tola/middleware"
)

func InitializeRoutes(router *gin.Engine) {
    router.GET("/ping", controllers.Ping)
    router.POST("/register", controllers.Register)
    router.POST("/login", controllers.Login)
    router.POST("/categories", controllers.AddCategory)
    router.GET("/categories", controllers.ListCategories)
    router.POST("/users/:id/categories", controllers.UpdateUserCategories)
    router.POST("/questions", controllers.AskQuestion)
    router.GET("/questions/:questionId/answers", controllers.ListAnswers)
    router.POST("/questions/:questionId/answers", controllers.AnswerQuestion)

    userGroup := router.Group("/user")
    userGroup.Use(middleware.AuthRequired) // Middleware pour vérifier l'authentification
    {
        userGroup.GET("/info", controllers.GetUserInfo)
        userGroup.GET("/posts", controllers.GetUserPosts)
    }
}
