package http

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/arrase21/crm/internal/domain"
	"github.com/arrase21/crm/internal/service"
	"github.com/gin-gonic/gin"
)

type EmployeeContractHandler struct {
	svc *service.EmployeeContractService
}

func NewEmployeeContractHandler(svc *service.EmployeeContractService) *EmployeeContractHandler {
	return &EmployeeContractHandler{svc}
}

type CreateEmployeeContractRequest struct {
	EmployeeID       uint   `json:"employee_id" binding:"required"`
	ContractTypeID   uint   `json:"contract_type_id"`
	BaseSalary       int64  `json:"base_salary" binding:"required,gt=0"`
	Currency         string `json:"currency" binding:"required,len=3"`
	CountryCode      string `json:"country_code" binding:"required,len=2"`
	StartDate        string `json:"start_date" binding:"required"`
	EndDate          string `json:"end_date"`
	WorkHoursPerDay  float64 `json:"work_hours_per_day" binding:"required,gt=0,lte=24"`
	WorkDaysPerWeek  float64 `json:"work_days_per_week" binding:"required,gt=0,lte=7"`
}

func (h *EmployeeContractHandler) Create(c *gin.Context) {
	var req CreateEmployeeContractRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format, use YYYY-MM-DD"})
		return
	}

	ec := &domain.EmployeeContract{
		EmployeeID:      req.EmployeeID,
		ContractTypeID:  req.ContractTypeID,
		BaseSalary:      req.BaseSalary,
		Currency:        req.Currency,
		CountryCode:     req.CountryCode,
		StartDate:       startDate,
		WorkHoursPerDay: req.WorkHoursPerDay,
		WorkDaysPerWeek: req.WorkDaysPerWeek,
		IsActive:        true,
	}

	if req.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format, use YYYY-MM-DD"})
			return
		}
		ec.EndDate = &endDate
	}

	if err := h.svc.Create(c.Request.Context(), ec); err != nil {
		if errors.Is(err, domain.ErrEmployeeNotFound) || errors.Is(err, domain.ErrContractTypeNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "contract created", "id": ec.ID})
}

func (h *EmployeeContractHandler) GetByID(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	contract, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrContractNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "contract not found"})
			return
		}
		internalError(c, err)
		return
	}

	c.JSON(http.StatusOK, contract)
}

func (h *EmployeeContractHandler) GetByEmployeeID(c *gin.Context) {
	empIDStr := c.Query("employee_id")
	if empIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "employee_id query parameter is required"})
		return
	}

	empID, err := strconv.ParseUint(empIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee_id"})
		return
	}

	contracts, err := h.svc.GetByEmployeeID(c.Request.Context(), uint(empID))
	if err != nil {
		internalError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"contracts": contracts})
}

func (h *EmployeeContractHandler) List(c *gin.Context) {
	page, limit := parsePagination(c)

	contracts, total, err := h.svc.List(c.Request.Context(), page, limit)
	if err != nil {
		internalError(c, err)
		return
	}

	respondPaginated(c, contracts, total, page, limit, "contracts")
}

type UpdateEmployeeContractRequest struct {
	ContractTypeID  *uint    `json:"contract_type_id"`
	BaseSalary      *int64   `json:"base_salary" binding:"omitempty,gt=0"`
	Currency        *string  `json:"currency" binding:"omitempty,len=3"`
	CountryCode     *string  `json:"country_code" binding:"omitempty,len=2"`
	StartDate       *string  `json:"start_date"`
	EndDate         *string  `json:"end_date"`
	IsActive        *bool    `json:"is_active"`
	WorkHoursPerDay *float64 `json:"work_hours_per_day" binding:"omitempty,gt=0,lte=24"`
	WorkDaysPerWeek *float64 `json:"work_days_per_week" binding:"omitempty,gt=0,lte=7"`
}

func (h *EmployeeContractHandler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var req UpdateEmployeeContractRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existing, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrContractNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "contract not found"})
			return
		}
		internalError(c, err)
		return
	}

	ec := &domain.EmployeeContract{
		ID:        existing.ID,
		TenantID:  existing.TenantID,
		CreatedAt: existing.CreatedAt,
	}

	if req.ContractTypeID != nil {
		ec.ContractTypeID = *req.ContractTypeID
	} else {
		ec.ContractTypeID = existing.ContractTypeID
	}
	if req.BaseSalary != nil {
		ec.BaseSalary = *req.BaseSalary
	} else {
		ec.BaseSalary = existing.BaseSalary
	}
	if req.Currency != nil {
		ec.Currency = *req.Currency
	} else {
		ec.Currency = existing.Currency
	}
	if req.CountryCode != nil {
		ec.CountryCode = *req.CountryCode
	} else {
		ec.CountryCode = existing.CountryCode
	}
	if req.WorkHoursPerDay != nil {
		ec.WorkHoursPerDay = *req.WorkHoursPerDay
	} else {
		ec.WorkHoursPerDay = existing.WorkHoursPerDay
	}
	if req.WorkDaysPerWeek != nil {
		ec.WorkDaysPerWeek = *req.WorkDaysPerWeek
	} else {
		ec.WorkDaysPerWeek = existing.WorkDaysPerWeek
	}
	if req.IsActive != nil {
		ec.IsActive = *req.IsActive
	} else {
		ec.IsActive = existing.IsActive
	}
	if req.StartDate != nil {
		d, err := time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format, use YYYY-MM-DD"})
			return
		}
		ec.StartDate = d
	} else {
		ec.StartDate = existing.StartDate
	}
	if req.EndDate != nil {
		if *req.EndDate == "" {
			ec.EndDate = nil
		} else {
			d, err := time.Parse("2006-01-02", *req.EndDate)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format, use YYYY-MM-DD"})
				return
			}
			ec.EndDate = &d
		}
	} else {
		ec.EndDate = existing.EndDate
	}

	if err := h.svc.Update(c.Request.Context(), ec); err != nil {
		if errors.Is(err, domain.ErrContractNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "contract not found"})
			return
		}
		internalError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "contract updated"})
}

func (h *EmployeeContractHandler) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, domain.ErrContractNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "contract not found"})
			return
		}
		internalError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
