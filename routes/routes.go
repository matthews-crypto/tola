package routes

import (
	"tola/controllers"
	"tola/middleware"

	"github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine) {
	router.GET("/ping", controllers.Ping)
	router.POST("/register", controllers.Register)
	router.POST("/login", controllers.Login)
	router.POST("/logout", controllers.Logout)
	router.POST("/categories", controllers.AddCategory)
	router.GET("/categories", controllers.ListCategories)
	router.POST("/users/:id/categories", controllers.UpdateUserCategories)
	router.POST("/questions", controllers.AskQuestion)
	router.GET("/questions/:id", controllers.GetQuestionDetails)
	router.GET("/questions/:id/answers", controllers.ListAnswers) // Modification ici
	router.GET("/search", controllers.Search)
	router.GET("/user/profile/:id", controllers.GetUserProfile)

	userGroup := router.Group("/user")
	userGroup.Use(middleware.AuthRequired) // Middleware pour vérifier l'authentification
	{
		userGroup.GET("/info", controllers.GetUserInfo)
		userGroup.GET("/posts", controllers.GetUserPosts)
		userGroup.POST("/update", controllers.UpdateUserProfile)
		userGroup.POST("/publications", controllers.CreatePublication)       // Nouvelle route pour créer une publication
		userGroup.GET("/publications", controllers.ListUserPublications)     // Nouvelle route pour lister les publications de l'utilisateur
		userGroup.DELETE("/publications/:id", controllers.DeletePublication) // Nouvelle route pour supprimer une publication
	}
}
