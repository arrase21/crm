package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/arrase21/crm/internal/domain"
	"github.com/arrase21/crm/internal/service"
	"github.com/gin-gonic/gin"
)

type EmployeeHandler struct {
	svc *service.EmployeeService
}

func NewEmployeeHandler(svc *service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{svc}
}

type CreateEmployeeRequest struct {
	UserID       uint  `json:"user_id" binding:"required"`
	DepartmentID uint  `json:"department_id" binding:"required"`
	PositionID   uint  `json:"position_id" binding:"required"`
	SupervisorID *uint `json:"supervisor_id"`
	IsActive     *bool `json:"is_active"`
}

func (h *EmployeeHandler) Create(c *gin.Context) {
	var req CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	emp := &domain.Employee{
		UserID:       req.UserID,
		DepartmentID: req.DepartmentID,
		PositionID:   req.PositionID,
		SupervisorID: req.SupervisorID,
	}
	if req.IsActive != nil {
		emp.IsActive = *req.IsActive
	}

	if err := h.svc.Create(c.Request.Context(), emp); err != nil {
		if errors.Is(err, domain.ErrEmployeeAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, domain.ErrUserNotFound) ||
			errors.Is(err, domain.ErrDepartmentNotFound) ||
			errors.Is(err, domain.ErrPositionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "employee created"})
}

func (h *EmployeeHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	emp, err := h.svc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrEmployeeNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "employee not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, emp)
}

func (h *EmployeeHandler) GetByUserID(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id query parameter is required"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	emp, err := h.svc.GetByUserID(c.Request.Context(), uint(userID))
	if err != nil {
		if errors.Is(err, domain.ErrEmployeeNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "employee not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, emp)
}

func (h *EmployeeHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	activeOnly := c.Query("active")
	supervisorID, _ := strconv.ParseUint(c.Query("supervisor_id"), 10, 32)

	var employees []domain.Employee
	var total int64
	var err error

	switch {
	case supervisorID > 0:
		employees, total, err = h.svc.ListBySupervisor(c.Request.Context(), uint(supervisorID), page, limit)
	case activeOnly == "true":
		employees, total, err = h.svc.ListActive(c.Request.Context(), page, limit)
	default:
		employees, total, err = h.svc.List(c.Request.Context(), page, limit)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, gin.H{
		"employees": employees,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

type UpdateEmployeeRequest struct {
	UserID       *uint  `json:"user_id"`
	DepartmentID *uint  `json:"department_id"`
	PositionID   *uint  `json:"position_id"`
	SupervisorID *uint  `json:"supervisor_id"`
	IsActive     *bool  `json:"is_active"`
}

func (h *EmployeeHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req UpdateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingEmp, err := h.svc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrEmployeeNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "employee not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	emp := &domain.Employee{
		ID:        existingEmp.ID,
		TenantID:  existingEmp.TenantID,
		CreatedAt: existingEmp.CreatedAt,
	}

	if req.UserID != nil {
		emp.UserID = *req.UserID
	} else {
		emp.UserID = existingEmp.UserID
	}
	if req.DepartmentID != nil {
		emp.DepartmentID = *req.DepartmentID
	} else {
		emp.DepartmentID = existingEmp.DepartmentID
	}
	if req.PositionID != nil {
		emp.PositionID = *req.PositionID
	} else {
		emp.PositionID = existingEmp.PositionID
	}
	if req.SupervisorID != nil {
		emp.SupervisorID = req.SupervisorID
	} else {
		emp.SupervisorID = existingEmp.SupervisorID
	}
	if req.IsActive != nil {
		emp.IsActive = *req.IsActive
	} else {
		emp.IsActive = existingEmp.IsActive
	}

	if err := h.svc.Update(c.Request.Context(), emp); err != nil {
		if errors.Is(err, domain.ErrEmployeeAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, domain.ErrUserNotFound) ||
			errors.Is(err, domain.ErrDepartmentNotFound) ||
			errors.Is(err, domain.ErrPositionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "employee updated"})
}

func (h *EmployeeHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.svc.Delete(c.Request.Context(), uint(id)); err != nil {
		if errors.Is(err, domain.ErrEmployeeNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "employee not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
