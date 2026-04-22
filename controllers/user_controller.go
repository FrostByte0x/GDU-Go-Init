package controllers

import (
	"errors"
	"net/http"
	"partage-projets/config"
	"partage-projets/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	var existingUser models.User
	result := config.DB.Where("email = ?", user.Email).First(&existingUser)
	if result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user email is already taken"})
		return
	}
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// something else went wrong
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}
	// update the password with the hashed values
	user.Password = string(hashedPassword)
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Success:": "User created successfully"})
}
