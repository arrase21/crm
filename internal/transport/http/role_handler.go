package http

import (
	"net/http"

	"github.com/arrase21/crm/internal/domain"
	"github.com/arrase21/crm/internal/service"
	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	svc *service.RoleService
}

func NewRoleHandler(svc *service.RoleService) *RoleHandler {
	return &RoleHandler{svc: svc}
}

type AssignRoleRequest struct {
	UserID   uint   `json:"user_id" binding:"required"`
	RoleName string `json:"role_name" binding:"required"`
}

func (h *RoleHandler) Assign(c *gin.Context) {
	var req AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.Assign(c.Request.Context(), req.UserID, req.RoleName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "role assigned"})
}

func (h *RoleHandler) Unassign(c *gin.Context) {
	var req AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.Unassign(c.Request.Context(), req.UserID, req.RoleName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "role unassigned"})
}

func (h *RoleHandler) ListRoles(c *gin.Context) {
	roles, err := h.svc.ListRoles(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, roles)
}

func (h *RoleHandler) UserRoles(c *gin.Context) {
	claims, ok := c.Request.Context().Value(domain.ClaimsKey).(*domain.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	roles, err := h.svc.GetUserRoles(c.Request.Context(), claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"roles": roles})
}
