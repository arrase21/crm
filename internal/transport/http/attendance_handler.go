package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/arrase21/crm/internal/domain"
	"github.com/arrase21/crm/internal/service"
	"github.com/gin-gonic/gin"
)

type AttendanceHandler struct {
	svc *service.AttendanceService
}

func NewAttendanceHandler(svc *service.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{svc: svc}
}

type CreateAttendanceRequest struct {
	EmployeeID uint   `json:"employee_id" binding:"required"`
	Date       string `json:"date" binding:"required"`
	CheckIn    string `json:"check_in"`
	CheckOut   string `json:"check_out"`
}

func (h *AttendanceHandler) Create(c *gin.Context) {
	var req CreateAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format (expected YYYY-MM-DD)"})
		return
	}

	a := &domain.Attendance{
		EmployeeID: req.EmployeeID,
		Date:       date,
	}

	if req.CheckIn != "" {
		t, err := time.Parse("15:04", req.CheckIn)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid check_in format (expected HH:MM)"})
			return
		}
		a.CheckIn = &t
	}
	if req.CheckOut != "" {
		t, err := time.Parse("15:04", req.CheckOut)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid check_out format (expected HH:MM)"})
			return
		}
		a.CheckOut = &t
	}

	if err := h.svc.Create(c.Request.Context(), a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, a)
}

func (h *AttendanceHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	a, err := h.svc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "attendance not found"})
		return
	}
	c.JSON(http.StatusOK, a)
}

func (h *AttendanceHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	employeeID, _ := strconv.ParseUint(c.Query("employee_id"), 10, 32)
	deptID, _ := strconv.ParseUint(c.Query("department_id"), 10, 32)

	var records []domain.Attendance
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, gin.H{
		"attendance": records,
		"pagination": gin.H{
			"page": page, "limit": limit,
			"total": total, "total_pages": totalPages,
		},
	})
}

type UpdateAttendanceRequest struct {
	Date     *string `json:"date"`
	CheckIn  *string `json:"check_in"`
	CheckOut *string `json:"check_out"`
}

func (h *AttendanceHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	existing, err := h.svc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "attendance not found"})
		return
	}

	var req UpdateAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Date != nil {
		t, err := time.Parse("2006-01-02", *req.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format"})
			return
		}
		existing.Date = t
	}
	if req.CheckIn != nil {
		t, err := time.Parse("15:04", *req.CheckIn)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid check_in format"})
			return
		}
		existing.CheckIn = &t
	}
	if req.CheckOut != nil {
		t, err := time.Parse("15:04", *req.CheckOut)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid check_out format"})
			return
		}
		existing.CheckOut = &t
	}

	if err := h.svc.Update(c.Request.Context(), existing); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, existing)
}

func (h *AttendanceHandler) Delete(c *gin.Context) {
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
