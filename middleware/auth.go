package middleware

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/sessions"
)

func AuthRequired(c *gin.Context) {
    session := sessions.Default(c)
    userID := session.Get("userID")
    if userID == nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non connecté"})
        c.Abort()
        return
    }
    c.Next()
}
