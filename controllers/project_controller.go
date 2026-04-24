package controllers

import (
	"errors"
	"log/slog"
	"net/http"
	"partage-projets/config"
	"partage-projets/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func GetProjects(c *gin.Context) {
	var projects []models.Project

	if err := config.DB.Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error: ": "could not load projects"})
		return
	}
	c.JSON(http.StatusOK, projects)
}

func PostProject(c *gin.Context) {
	var project models.Project

	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	if err := config.DB.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		slog.Warn(err.Error())
		return
	}

	c.JSON(http.StatusOK, project)
}

func GetProject(c *gin.Context) {
	var project models.Project

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	if err := config.DB.First(&project, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Project could not be found"})
		return
	}
	c.JSON(http.StatusOK, project)
}

func PutProject(c *gin.Context) {
	var project models.Project
	// Request will be received on /projects/id
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	if err := config.DB.First(&project, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//if not found, return it
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}
		// handle other errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Project could not be found"})
		return
	}
	var input models.ProjectUpdateInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	updates := make(map[string]any)

	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.Description != nil {
		updates["description"] = *input.Description
	}
	if input.Image != nil {
		updates["image"] = *input.Image
	}
	if input.Skills != nil {
		updates["skills"] = datatypes.JSONSlice[string](*input.Skills)
	}
	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "there are no fields to update"})
		return
	}
	if err := config.DB.Model(&project).Updates(updates).Error; err != nil {
		// print err to client for debug
		c.JSON(http.StatusInternalServerError, gin.H{"error:": err.Error()})
		return
	}
	c.JSON(http.StatusOK, project)
}

func DeleteProject(c *gin.Context) {
	var project models.Project // this var is empty, the request param ID will be used to retrieve it in the DB and
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
	}
	if err := config.DB.First(&project, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//if not found, return it
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}
		// handle other errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Project could not be found"})
		return
	}

	if err := config.DB.Delete(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error: ": "error deleting project"})
		return
	}
	c.JSON(http.StatusOK, project)
}
