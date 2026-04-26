package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/arrase21/crm/internal/domain"
	"github.com/arrase21/crm/internal/service"
	"github.com/gin-gonic/gin"
)

type DepartmentHandler struct {
	svc *service.DepartmentService
}

func NewDepartmentHandler(svc *service.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{svc}
}

type CreateDepartmentRequest struct {
	Name     string `json:"name" binding:"required,max=100"`
	Code     string `json:"code" binding:"required,max=100"`
	IsActive *bool  `json:"is_active"`
}

func (h *DepartmentHandler) Create(c *gin.Context) {
	var req CreateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	dept := &domain.Department{
		Name: req.Name,
		Code: req.Code,
	}
	if req.IsActive != nil {
		dept.IsActive = *req.IsActive
	}
	if err := h.svc.Create(c.Request.Context(), dept); err != nil {
		if errors.Is(err, domain.ErrDepartmentCodeExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "department create"})
}

func (h *DepartmentHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	dept, err := h.svc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrDepartmentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "department not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dept)
}

func (h *DepartmentHandler) GetByCode(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code query parameter is required"})
		return
	}
	dept, err := h.svc.GetByCode(c.Request.Context(), code)
	if err != nil {
		if errors.Is(err, domain.ErrDepartmentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "department code not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dept)
}

func (h *DepartmentHandler) GetByName(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name query parameter is required"})
		return
	}
	dept, err := h.svc.GetByName(c.Request.Context(), name)
	if err != nil {
		if errors.Is(err, domain.ErrDepartmentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "department name not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dept)
}

func (h *DepartmentHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	dept, total, err := h.svc.List(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}
	c.JSON(http.StatusOK, gin.H{
		"departments": dept,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

type UpdateDepartmentRequest struct {
	Name     *string `json:"name" binding:"omitempty,max=100"`
	Code     *string `json:"code" binding:"omitempty,max=20"`
	IsActive *bool   `json:"is_active"`
}

func (h *DepartmentHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req UpdateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	existingDept, err := h.svc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrDepartmentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "department  not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	dept := &domain.Department{
		ID:        existingDept.ID,
		TenantID:  existingDept.TenantID,
		CreatedAt: existingDept.CreatedAt,
	}
	if req.Name != nil {
		dept.Name = *req.Name
	} else {
		dept.Name = existingDept.Name
	}
	if req.Code != nil {
		dept.Code = *req.Code
	} else {
		dept.Code = existingDept.Code
	}
	if req.IsActive != nil {
		dept.IsActive = *req.IsActive
	} else {
		dept.IsActive = existingDept.IsActive
	}
	if err := h.svc.Update(c.Request.Context(), dept); err != nil {
		if errors.Is(err, domain.ErrDepartmentCodeExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "department updated"})
}

func (h *DepartmentHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}
	if err := h.svc.Delete(c.Request.Context(), uint(id)); err != nil {
		if errors.Is(err, domain.ErrDepartmentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "department not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
