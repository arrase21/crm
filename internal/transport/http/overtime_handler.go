package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/arrase21/crm/internal/domain"
	"github.com/arrase21/crm/internal/service"
	"github.com/gin-gonic/gin"
)

type OvertimeHandler struct {
	svc *service.OvertimeService
}

func NewOvertimeHandler(svc *service.OvertimeService) *OvertimeHandler {
	return &OvertimeHandler{svc: svc}
}

type CreateOvertimeRequest struct {
	EmployeeID uint    `json:"employee_id" binding:"required"`
	Date       string  `json:"date" binding:"required"`
	Hours      float64 `json:"hours" binding:"required,min=0.5"`
	Reason     string  `json:"reason"`
}

func (h *OvertimeHandler) Create(c *gin.Context) {
	var req CreateOvertimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format (expected YYYY-MM-DD)"})
		return
	}

	o := &domain.Overtime{
		EmployeeID: req.EmployeeID,
		Date:       date,
		Hours:      req.Hours,
		Reason:     req.Reason,
	}

	if err := h.svc.Create(c.Request.Context(), o); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, o)
}

func (h *OvertimeHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	o, err := h.svc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "overtime not found"})
		return
	}
	c.JSON(http.StatusOK, o)
}

func (h *OvertimeHandler) List(c *gin.Context) {
	page, limit := parsePagination(c)
	employeeID, _ := strconv.ParseUint(c.Query("employee_id"), 10, 32)
	deptID, _ := strconv.ParseUint(c.Query("department_id"), 10, 32)

	var records []domain.Overtime
	var total int64
	var err error

	switch {
	case employeeID > 0:
		records, total, err = h.svc.ListByEmployee(c.Request.Context(), uint(employeeID), page, limit)
	case deptID > 0:
		records, total, err = h.svc.ListByDepartment(c.Request.Context(), uint(deptID), page, limit)
	default:
		records, total, err = h.svc.List(c.Request.Context(), page, limit)
	}

	if err != nil {
		internalError(c, err)
		return
	}

	respondPaginated(c, records, total, page, limit, "overtime")
}

type UpdateOvertimeRequest struct {
	Hours  *float64 `json:"hours"`
	Reason *string  `json:"reason"`
}

func (h *OvertimeHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	existing, err := h.svc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "overtime not found"})
		return
	}

	var req UpdateOvertimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Hours != nil {
		existing.Hours = *req.Hours
	}
	if req.Reason != nil {
		existing.Reason = *req.Reason
	}

	if err := h.svc.Update(c.Request.Context(), existing); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, existing)
}

func (h *OvertimeHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.svc.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
