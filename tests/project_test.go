package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"partage-projets/config"
	"partage-projets/controllers"
	"partage-projets/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite" // Since I'm running Go on Windows
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const projectName string = "Project test"
const projectDescription string = "Projet Description"
const projectComment string = "Comment test"

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		fmt.Printf("Failed to open database: %v", err)
		return nil
	}
	db.AutoMigrate(&models.Project{}, &models.Comment{})

	project := models.Project{Name: projectName, Description: projectDescription}
	db.Create(&project)

	comment := models.Comment{ProjectId: project.ID, Content: projectComment}
	db.Create(&comment)

	return db
}

func TestGetProject(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config.DB = setupTestDB()

	r := gin.Default()
	r.GET("/projects", controllers.GetProjects)

	req, _ := http.NewRequest(http.MethodGet, "/projects", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	body := w.Body.String()
	assert.Contains(t, body, projectName)
	assert.Contains(t, body, projectDescription)
}

func TestPostProject(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config.DB = setupTestDB()

	r := gin.Default()
	r.POST("/projects", controllers.PostProject)

	project := map[string]any{
		"name":        projectName,
		"description": projectDescription,
		"skills":      []string{"Go", "Testify", "Testing"},
	}

	data, err := json.Marshal(project)
	if err != nil {
		slog.Warn(err.Error())
	}
	req, _ := http.NewRequest(http.MethodPost, "/projects", bytes.NewBuffer(data))
	req.Header.Set("content-type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), projectName)
	assert.Contains(t, w.Body.String(), projectDescription)
	assert.Contains(t, w.Body.String(), "Testify")
}
