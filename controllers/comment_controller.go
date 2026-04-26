package controllers

import (
	"net/http"
	"partage-projets/config"
	"partage-projets/models"

	"github.com/gin-gonic/gin"
)

func PostComment(c *gin.Context) {
	var comment models.Comment

	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}
	// type assertion!

	userIdInt, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	comment.UserID = uint(userIdInt)

	if err := config.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error saving comment"})
		return
	}
	c.JSON(http.StatusCreated, comment)
}
