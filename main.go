package main

import (
    "tola/routes"
    "tola/utils"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/memstore"
)

func main() {
    utils.ConnectDB()
    utils.InitCollections() // Initialiser les collections

    router := gin.Default()

    // Configurer le middleware de sessions
    store := memstore.NewStore([]byte("secret"))
    router.Use(sessions.Sessions("mysession", store))

    // Servir les fichiers statiques
    router.Static("/public", "./public")

    // Route pour la page d'accueil
    router.GET("/", func(c *gin.Context) {
        c.File("./public/index.html")
    })

    // Route pour ajouter une catégorie
    router.GET("/add_category", func(c *gin.Context) {
        c.File("./public/add_category.html")
    })

    // Route pour sélectionner des catégories après inscription
    router.GET("/select_categories", func(c *gin.Context) {
        c.File("./public/select_categories.html")
    })

    // Route pour poser une question
    router.GET("/ask_question", func(c *gin.Context) {
        c.File("./public/ask_question.html")
    })

    // Route pour afficher les questions
    router.GET("/questions", func(c *gin.Context) {
        c.File("./public/questions.html")
    })

    // Route pour répondre à une question
    router.GET("/answer_question", func(c *gin.Context) {
        c.File("./public/answer_question.html")
    })

    // Route pour la page de profil
    router.GET("/profile", func(c *gin.Context) {
        c.File("./public/profile.html")
    })

    routes.InitializeRoutes(router)
    router.Run(":8080")
}
