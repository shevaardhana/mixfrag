package controllers

import (
	"database/sql"
	"net/http"

	"mixfrag/database"

	"github.com/gin-gonic/gin"
)

// Middleware Basic Auth
func BasicAuthMiddleware(c *gin.Context) {
	username, password, ok := c.Request.BasicAuth()
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		c.Abort()
		return
	}

	var hashedPassword string
	// ⚠️ pakai schema public biar aman
	err := database.DB.QueryRow("SELECT password FROM public.users WHERE username=$1", username).Scan(&hashedPassword)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		c.Abort()
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	// cek password
	if hashedPassword != password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		c.Abort()
		return
	}

	// simpan username di context
	c.Set("username", username)
	c.Next()
}

// Contoh endpoint protected
func Profile(c *gin.Context) {
	username := c.GetString("username")
	c.JSON(http.StatusOK, gin.H{
		"message":  "Welcome to your profile!",
		"username": username,
	})
}
