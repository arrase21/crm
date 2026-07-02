package http

import (
	"errors"
	"net/http"
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
	Password  string `json:"password" binding:"required,min=4"`
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

	user := &domain.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Dni:       req.Dni,
		Gender:    req.Gender,
		Phone:     req.Phone,
		Email:     req.Email,
		BirthDay:  birth,
	}

	if err := h.svc.Create(c.Request.Context(), user, req.Password); err != nil {
		if errors.Is(err, domain.ErrDniAlreadyExist) ||
			errors.Is(err, domain.ErrEmailAlreadyExist) ||
			errors.Is(err, domain.ErrPhoneAlreadyExist) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user created"})
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	usr, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		internalError(c, err)
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
		internalError(c, err)
		return
	}

	c.JSON(http.StatusOK, usr)
}

func (h *UserHandler) List(c *gin.Context) {
	page, limit := parsePagination(c)

	usrs, total, err := h.svc.List(c.Request.Context(), page, limit)
	if err != nil {
		internalError(c, err)
		return
	}

	respondPaginated(c, usrs, total, page, limit, "users")
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
	id, ok := parseID(c)
	if !ok {
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingUser, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		internalError(c, err)
		return
	}

	user := &domain.User{
		ID:        existingUser.ID,
		TenantID:  existingUser.TenantID,
		CreatedAt: existingUser.CreatedAt,
	}

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

	if err := h.svc.Update(c.Request.Context(), user); err != nil {
		if errors.Is(err, domain.ErrDniAlreadyExist) ||
			errors.Is(err, domain.ErrEmailAlreadyExist) ||
			errors.Is(err, domain.ErrPhoneAlreadyExist) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		internalError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user updated"})
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		internalError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
