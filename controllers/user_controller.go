package controllers

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"partage-projets/config"
	"partage-projets/models"
	"partage-projets/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type CustomClaim struct {
	UserID uint
	jwt.RegisteredClaims
}

func Login(c *gin.Context) {
	// receive username and password, in params?
	// We check the user exists and the email matches
	// we hash the password and check that it results in the same one way hash we have in the database
	// If all these steps succeed, we return a JWT to the user that they can use to authenticate.

	var user models.User
	// unmarshal the json into user
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	// Check that username + hashed password matches the entry in the database
	// this one is harder, we need to check for 2 fields, not just the primary key which is often the Id
	// Return
	var existingUser models.User
	if err := config.DB.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email or password is not valid"})
		return
	}
	// Validate the input password with the one already hashed.
	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email or password is not valid"})
		return
	}
	claim := &CustomClaim{
		UserID: existingUser.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	// this does not look safe
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		slog.Warn(err.Error())
		return
	}
	c.JSON(http.StatusOK, tokenString)
}

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

	if err := utils.Validatepassword(user.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
