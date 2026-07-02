package http

import (
	"errors"
	"net/http"

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
	Code     string `json:"code" binding:"required,max=20"`
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
	c.JSON(http.StatusCreated, gin.H{"message": "department created"})
}

func (h *DepartmentHandler) GetByID(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	dept, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrDepartmentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "department not found"})
			return
		}
		internalError(c, err)
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
		internalError(c, err)
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
		internalError(c, err)
		return
	}
	c.JSON(http.StatusOK, dept)
}

func (h *DepartmentHandler) List(c *gin.Context) {
	page, limit := parsePagination(c)

	dept, total, err := h.svc.List(c.Request.Context(), page, limit)
	if err != nil {
		internalError(c, err)
		return
	}

	respondPaginated(c, dept, total, page, limit, "departments")
}

type UpdateDepartmentRequest struct {
	Name     *string `json:"name" binding:"omitempty,max=100"`
	Code     *string `json:"code" binding:"omitempty,max=20"`
	IsActive *bool   `json:"is_active"`
}

func (h *DepartmentHandler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req UpdateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	existingDept, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrDepartmentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "department not found"})
			return
		}
		internalError(c, err)
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
		internalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "department updated"})
}

func (h *DepartmentHandler) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, domain.ErrDepartmentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "department not found"})
			return
		}
		internalError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
