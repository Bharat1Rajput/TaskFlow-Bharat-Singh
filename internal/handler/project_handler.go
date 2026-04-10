package handler

import (
	"github.com/Bharat1Rajput/taskflow-backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProjectHandler struct {
	service *service.ProjectService
}

func NewProjectHandler(s *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: s}
}

func (h *ProjectHandler) Create(c *gin.Context) {
	var req struct {
		Name        string  `json:"name"`
		Description *string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.Name == "" {
		c.JSON(400, gin.H{
			"error":  "validation failed",
			"fields": gin.H{"name": "required"},
		})
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	err := h.service.Create(c.Request.Context(), req.Name, req.Description, userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "project created"})
}

func (h *ProjectHandler) List(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	projects, err := h.service.List(c.Request.Context(), userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, projects)
}

func (h *ProjectHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := uuid.Parse(idParam)

	var req struct {
		Name        string  `json:"name"`
		Description *string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "validation failed"})
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	err := h.service.Update(c.Request.Context(), id, req.Name, req.Description, userID)
	if err != nil {
		if err == service.ErrForbidden {
			c.JSON(403, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "updated"})
}

func (h *ProjectHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := uuid.Parse(idParam)

	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	err := h.service.Delete(c.Request.Context(), id, userID)
	if err != nil {
		if err == service.ErrForbidden {
			c.JSON(403, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "deleted"})
}
