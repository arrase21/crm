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

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc}
}

// parseBirthDay intenta parsear una fecha con múltiples formatos
func parseBirthDay(dateStr string) (time.Time, error) {
	birth, err := time.Parse(time.RFC3339Nano, dateStr)
	if err == nil {
		return birth, nil
	}
	birth, err = time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, errors.New("invalid birth_day format (expected RFC3339 or YYYY-MM-DD)")
	}
	return birth, nil
}

type CreateUserRequest struct {
	FirstName string `json:"first_name" binding:"required,max=30"`
	LastName  string `json:"last_name" binding:"required,max=40"`
	Dni       string `json:"dni" binding:"required,max=20"`
	Gender    string `json:"gender" binding:"required,oneof=M F"`
	Phone     string `json:"phone" binding:"required,max=15"`
	Email     string `json:"email" binding:"required,email,max=50"`
	BirthDay  string `json:"birth_day" binding:"required"`
}

func (h *UserHandler) Create(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	birth, err := parseBirthDay(req.BirthDay)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Construir objeto de dominio
	user := &domain.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Dni:       req.Dni,
		Gender:    req.Gender,
		Phone:     req.Phone,
		Email:     req.Email,
		BirthDay:  birth,
	}

	user.Normalize()

	if err := user.ValidateAll(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.Create(c.Request.Context(), user); err != nil {
		if errors.Is(err, domain.ErrDniAlreadyExist) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user created"})
}

func (h *UserHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	usr, err := h.svc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, usr)
}

func (h *UserHandler) GetByDni(c *gin.Context) {
	dni := c.Query("dni")
	if dni == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dni query parameter is required"})
		return
	}

	usr, err := h.svc.GetByDNI(c.Request.Context(), dni)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, usr)
}

func (h *UserHandler) List(c *gin.Context) {
	// Parsear paginación
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	usrs, total, err := h.svc.List(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, gin.H{
		"users": usrs,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

type UpdateUserRequest struct {
	FirstName *string `json:"first_name" binding:"omitempty,max=30"`
	LastName  *string `json:"last_name" binding:"omitempty,max=40"`
	Dni       *string `json:"dni" binding:"omitempty,max=20"`
	Gender    *string `json:"gender" binding:"omitempty,oneof=M F"`
	Phone     *string `json:"phone" binding:"omitempty,max=15"`
	Email     *string `json:"email" binding:"omitempty,email,max=50"`
	BirthDay  *string `json:"birth_day"`
}

func (h *UserHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obtener usuario existente con sus roles
	existingUser, err := h.svc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Preservar campos que no deben cambiar
	// originalRoles := existingUser.Roles
	originalTenantID := existingUser.TenantID
	originalID := existingUser.ID
	originalCreatedAt := existingUser.CreatedAt

	// Construir usuario actualizado
	user := &domain.User{
		ID:        originalID,
		TenantID:  originalTenantID,
		CreatedAt: originalCreatedAt,
		// Roles:     originalRoles, // Preservar roles originales
	}

	// Actualizar solo los campos enviados
	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	} else {
		user.FirstName = existingUser.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	} else {
		user.LastName = existingUser.LastName
	}
	if req.Dni != nil {
		user.Dni = *req.Dni
	} else {
		user.Dni = existingUser.Dni
	}
	if req.Gender != nil {
		user.Gender = *req.Gender
	} else {
		user.Gender = existingUser.Gender
	}
	if req.Phone != nil {
		user.Phone = *req.Phone
	} else {
		user.Phone = existingUser.Phone
	}
	if req.Email != nil {
		user.Email = *req.Email
	} else {
		user.Email = existingUser.Email
	}
	if req.BirthDay != nil {
		birth, err := parseBirthDay(*req.BirthDay)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user.BirthDay = birth
	} else {
		user.BirthDay = existingUser.BirthDay
	}

	user.Normalize()
	if err := user.ValidateAll(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.Update(c.Request.Context(), user); err != nil {
		if errors.Is(err, domain.ErrDniAlreadyExist) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user updated"})
}

func (h *UserHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	if err := h.svc.Delete(c.Request.Context(), uint(id)); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
