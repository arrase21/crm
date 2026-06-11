package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/arrase21/crm/internal/domain"
	"github.com/arrase21/crm/internal/service"
	"github.com/gin-gonic/gin"
)

type PositionHandler struct {
	svc *service.PositionService
}

func NewPositionHandler(svc *service.PositionService) *PositionHandler {
	return &PositionHandler{svc}
}

type CreatePositionRequest struct {
	Name         string `json:"name" binding:"required,max=100"`
	Description  string `json:"description" binding:"required,max=255"`
	DepartmentID uint   `json:"department_id"`
	IsActive     *bool  `json:"is_active"`
}

func (h *PositionHandler) Create(c *gin.Context) {
	var req CreatePositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	position := &domain.Position{
		Name:         req.Name,
		Description:  req.Description,
		DepartmentID: req.DepartmentID,
	}
	if req.IsActive != nil {
		position.IsActive = *req.IsActive
	}
	if err := h.svc.Create(c.Request.Context(), position); err != nil {
		if errors.Is(err, domain.ErrPositionNameExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "position created"})
}

func (h *PositionHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// Check if we should include department in response
	withDepartment := c.Query("include_department")

	var position *domain.Position
	if withDepartment == "true" {
		position, err = h.svc.GetByIDWithDepartment(c.Request.Context(), uint(id))
	} else {
		position, err = h.svc.GetByID(c.Request.Context(), uint(id))
	}

	if err != nil {
		if errors.Is(err, domain.ErrPositionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "position not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, position)
}

func (h *PositionHandler) GetByName(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name query parameter is required"})
		return
	}
	position, err := h.svc.GetByName(c.Request.Context(), name)
	if err != nil {
		if errors.Is(err, domain.ErrPositionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "position name not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, position)
}

func (h *PositionHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	departmentIDStr := c.Query("department_id")
	var positions []domain.Position
	var total int64
	var err error

	if departmentIDStr != "" {
		departmentID, err := strconv.ParseUint(departmentIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid department_id"})
			return
		}
		positions, total, err = h.svc.ListByDepartment(c.Request.Context(), uint(departmentID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		positions, total, err = h.svc.List(c.Request.Context(), page, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}
	c.JSON(http.StatusOK, gin.H{
		"positions": positions,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

type UpdatePositionRequest struct {
	Name         *string `json:"name" binding:"omitempty,max=100"`
	Description  *string `json:"description" binding:"omitempty,max=255"`
	DepartmentID *uint   `json:"department_id"`
	IsActive     *bool   `json:"is_active"`
}

func (h *PositionHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req UpdatePositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	existingPosition, err := h.svc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrPositionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "position not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	position := &domain.Position{
		ID:        existingPosition.ID,
		TenantID:  existingPosition.TenantID,
		CreatedAt: existingPosition.CreatedAt,
	}
	if req.Name != nil {
		position.Name = *req.Name
	} else {
		position.Name = existingPosition.Name
	}
	if req.Description != nil {
		position.Description = *req.Description
	} else {
		position.Description = existingPosition.Description
	}
	if req.DepartmentID != nil {
		position.DepartmentID = *req.DepartmentID
	} else {
		position.DepartmentID = existingPosition.DepartmentID
	}
	if req.IsActive != nil {
		position.IsActive = *req.IsActive
	} else {
		position.IsActive = existingPosition.IsActive
	}
	if err := h.svc.Update(c.Request.Context(), position); err != nil {
		if errors.Is(err, domain.ErrPositionNameExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "position updated"})
}

func (h *PositionHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.svc.Delete(c.Request.Context(), uint(id)); err != nil {
		if errors.Is(err, domain.ErrPositionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "position not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
