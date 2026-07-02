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

type PayrollHandler struct {
	svc *service.PayrollService
}

func NewPayrollHandler(svc *service.PayrollService) *PayrollHandler {
	return &PayrollHandler{svc}
}

type CalculatePayrollRequest struct {
	ContractID  uint   `json:"contract_id" binding:"required"`
	PeriodStart string `json:"period_start" binding:"required"`
	PeriodEnd   string `json:"period_end" binding:"required"`
	Bonuses     int64  `json:"bonuses"`
	Commissions int64  `json:"commissions"`
}

func (h *PayrollHandler) CalculatePayroll(c *gin.Context) {
	var req CalculatePayrollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	periodStart, err := time.Parse("2006-01-02", req.PeriodStart)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid period_start format, use YYYY-MM-DD"})
		return
	}

	periodEnd, err := time.Parse("2006-01-02", req.PeriodEnd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid period_end format, use YYYY-MM-DD"})
		return
	}

	record, err := h.svc.CalculatePayroll(c.Request.Context(), req.ContractID, periodStart, periodEnd, req.Bonuses, req.Commissions)
	if err != nil {
		if errors.Is(err, domain.ErrContractNotFound) || errors.Is(err, domain.ErrCountryParamNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, domain.ErrCalculatorNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		internalError(c, err)
		return
	}

	c.JSON(http.StatusCreated, record)
}

func (h *PayrollHandler) GetPayrollRecords(c *gin.Context) {
	contractIDStr := c.Query("contract_id")
	if contractIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "contract_id query parameter is required"})
		return
	}

	contractID, err := strconv.ParseUint(contractIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid contract_id"})
		return
	}

	page, limit := parsePagination(c)

	records, total, err := h.svc.GetPayrollRecords(c.Request.Context(), uint(contractID), page, limit)
	if err != nil {
		internalError(c, err)
		return
	}

	respondPaginated(c, records, total, page, limit, "payroll_records")
}
