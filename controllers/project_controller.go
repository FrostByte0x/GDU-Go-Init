package controllers

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"partage-projets/config"
	"partage-projets/models"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// GetProjects godoc
// @Description Return a project
// @Tags Projects
// @Produce json
// @Success 200 {array} models.Project
// @Security BearerAuth
// @Router /projects [get]
func GetProjects(c *gin.Context) {
	var projects []models.Project

	if err := config.DB.Preload("Comments").Preload("Likes").Find(&projects).Error; err != nil {
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
	file, err := c.FormFile("image")
	if err == nil {
		path := "/uploads" + file.Filename
		if err := c.SaveUploadedFile(file, path); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot save image"})
			return
		}
		image, _ := imaging.Open(path)
		resized := imaging.Resize(image, 800, 0, imaging.Lanczos)
		if err := imaging.Save(resized, path); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot resize image"})
			slog.Warn(err.Error())
			return
		}
		project.Image = path
	}

	if err := config.DB.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		slog.Warn(err.Error())
		return
	}

	c.JSON(http.StatusCreated, project)
}

func GetProject(c *gin.Context) {
	var project models.Project

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	if err := config.DB.Preload("Comments").Preload("Likes").First(&project, id).Error; err != nil {
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

	if err := c.ShouldBind(&input); err != nil {
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
	// if input.Image != nil {
	// 	updates["image"] = *input.Image
	// }

	// handle image update
	file, err := c.FormFile("image")
	if err == nil {
		path := "/uploads" + file.Filename
		if err := c.SaveUploadedFile(file, path); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot save image"})
			return
		}
		image, _ := imaging.Open(path)
		resized := imaging.Resize(image, 800, 0, imaging.Lanczos)
		if err := imaging.Save(resized, path); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot resize image"})
			slog.Warn(err.Error())
			return
		}

		if project.Image != "" {
			_ = os.Remove(project.Image)
		}
		updates["Image"] = path
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
	var project models.Project
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

func LikeProject(c *gin.Context) {
	var project models.Project
	var user models.User
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
	}
	if err := config.DB.Preload("Likes").First(&project, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//if not found, return it
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}
		// handle other errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Project could not be found"})
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

	if err := config.DB.First(&user, uint(userIdInt)).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	liked := false
	for _, u := range project.Likes {
		if u.Id == user.Id {
			liked = true
			break
		}
	}
	if liked {
		if err := config.DB.Model(&project).Association("Likes").Delete(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error removing the like on this project"})
			slog.Warn(err.Error())
			return
		}
		c.JSON(http.StatusOK, "Sucess removing like on the project")
	} else {
		if err := config.DB.Model(&project).Association("Likes").Append(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error adding the like on this project"})
			slog.Warn(err.Error())
			return
		}
		c.JSON(http.StatusOK, "Sucess adding like on the project")
	}

}
