package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Bharat1Rajput/taskflow-backend/internal/model"
	"github.com/Bharat1Rajput/taskflow-backend/internal/service"
)

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(s *service.TaskService) *TaskHandler {
	return &TaskHandler{service: s}
}

type CreateTaskRequest struct {
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	AssigneeID  *uuid.UUID `json:"assignee_id"`
	DueDate     *time.Time `json:"due_date"`
}

type UpdateTaskRequest struct {
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	AssigneeID  *uuid.UUID `json:"assignee_id"`
	DueDate     *time.Time `json:"due_date"`
}

func (h *TaskHandler) Create(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid project id"})
		return
	}

	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Title == "" {
		c.JSON(400, gin.H{
			"error":  "validation failed",
			"fields": gin.H{"title": "required"},
		})
		return
	}

	task := model.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
		ProjectID:   projectID,
		AssigneeID:  req.AssigneeID,
		DueDate:     req.DueDate,
	}

	err = h.service.Create(c.Request.Context(), task)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "task created"})
}

func (h *TaskHandler) List(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid project id"})
		return
	}

	// Query params
	var status *string
	if s := c.Query("status"); s != "" {
		status = &s
	}

	var assignee *uuid.UUID
	if a := c.Query("assignee"); a != "" {
		uid, err := uuid.Parse(a)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid assignee id"})
			return
		}
		assignee = &uid
	}

	tasks, err := h.service.List(
		c.Request.Context(),
		projectID,
		status,
		assignee,
	)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, tasks)
}

func (h *TaskHandler) Update(c *gin.Context) {
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid task id"})
		return
	}

	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "validation failed"})
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	task := model.Task{
		ID:          taskID,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
		AssigneeID:  req.AssigneeID,
		DueDate:     req.DueDate,
	}

	err = h.service.Update(c.Request.Context(), task, userID)
	if err != nil {
		if err == service.ErrForbidden {
			c.JSON(403, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "task updated"})
}

func (h *TaskHandler) Delete(c *gin.Context) {
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid task id"})
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	err = h.service.Delete(c.Request.Context(), taskID, userID)
	if err != nil {
		if err == service.ErrForbidden {
			c.JSON(403, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "task deleted"})
}
