package users

import (
	"net/http"
	"strconv"

	"github.com/evencycu/TeamsNotifyGoV2/internal/api/services"
	"github.com/evencycu/TeamsNotifyGoV2/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles user-related HTTP requests
type Handler struct {
	userService services.UserService
}

// NewHandler creates a new user handler
func NewHandler(userService services.UserService) *Handler {
	return &Handler{
		userService: userService,
	}
}

// RegisterRoutes registers user routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/users")
	{
		users.POST("", h.CreateUser)
		users.GET("", h.ListUsers)
		users.GET("/:id", h.GetUser)
		users.PUT("/:id", h.UpdateUser)
		users.DELETE("/:id", h.DeleteUser)
		users.PATCH("/:id/password", h.ChangePassword)
		users.GET("/company/:companyId", h.GetUsersByCompany)
		users.GET("/role/:role", h.GetUsersByRole)
	}
}

// CreateUserRequest represents a create user request
type CreateUserRequest struct {
	CompanyID uuid.UUID `json:"companyId" validate:"required"`
	Email     string    `json:"email" validate:"required,email"`
	Name      string    `json:"name" validate:"required,min=2,max=255"`
	Role      string    `json:"role" validate:"required,oneof=admin manager user"`
	Password  string    `json:"password" validate:"required,min=8"`
}

// UpdateUserRequest represents an update user request
type UpdateUserRequest struct {
	Email string `json:"email" validate:"omitempty,email"`
	Name  string `json:"name" validate:"omitempty,min=2,max=255"`
	Role  string `json:"role" validate:"omitempty,oneof=admin manager user"`
}

// ChangePasswordRequest represents a change password request
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required,min=8"`
}

// CreateUser creates a new user
func (h *Handler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Convert request to user model
	user := &database.User{
		CompanyID: req.CompanyID,
		Email:     req.Email,
		Name:      req.Name,
		Role:      req.Role,
		Status:    "active", // Default status
		// Password will be hashed in service layer
	}

	// Create user
	createReq := &services.CreateRequest[database.User]{
		Data: *user,
	}

	createdUser, err := h.userService.Create(c.Request.Context(), createReq)
	if err != nil {
		if _, ok := err.(services.ConflictError); ok {
			c.JSON(http.StatusConflict, gin.H{
				"error": "User with this email already exists",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create user",
			"details": err.Error(),
		})
		return
	}

	// Remove sensitive fields from response
	createdUser.PasswordHash = nil

	c.JSON(http.StatusCreated, gin.H{
		"data":    createdUser,
		"message": "User created successfully",
	})
}

// ListUsers lists users with pagination
func (h *Handler) ListUsers(c *gin.Context) {
	// Parse pagination parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	sortBy := c.DefaultQuery("sort_by", "created_at")
	order := c.DefaultQuery("order", "desc")
	search := c.Query("search")
	status := c.Query("status")

	// Create list request
	req := &services.ListRequest{
		Limit:  limit,
		Offset: offset,
		SortBy: sortBy,
		Order:  order,
		Search: search,
		Status: status,
	}

	// Get users
	users, err := h.userService.List(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to list users",
			"details": err.Error(),
		})
		return
	}

	// Remove sensitive fields from response
	for _, user := range users {
		user.PasswordHash = nil
	}

	// Get total count
	countReq := &services.CountRequest{
		Search: search,
		Status: status,
	}
	total, err := h.userService.Count(c.Request.Context(), countReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to count users",
			"details": err.Error(),
		})
		return
	}

	// Calculate pages
	pages := int((total + int64(limit) - 1) / int64(limit))

	c.JSON(http.StatusOK, gin.H{
		"data": users,
		"pagination": gin.H{
			"total":  total,
			"limit":  limit,
			"offset": offset,
			"pages":  pages,
		},
	})
}

// GetUser gets a user by ID
func (h *Handler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	user, err := h.userService.GetByID(c.Request.Context(), id)
	if err != nil {
		if _, ok := err.(services.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get user",
			"details": err.Error(),
		})
		return
	}

	// Remove sensitive fields from response
	user.PasswordHash = nil

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

// UpdateUser updates a user
func (h *Handler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Get existing user
	existingUser, err := h.userService.GetByID(c.Request.Context(), id)
	if err != nil {
		if _, ok := err.(services.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get user",
			"details": err.Error(),
		})
		return
	}

	// Update fields
	if req.Email != "" {
		existingUser.Email = req.Email
	}
	if req.Name != "" {
		existingUser.Name = req.Name
	}
	if req.Role != "" {
		existingUser.Role = req.Role
	}

	// Update user
	updateReq := &services.UpdateRequest[database.User]{
		Data: *existingUser,
	}

	updatedUser, err := h.userService.Update(c.Request.Context(), id, updateReq)
	if err != nil {
		if _, ok := err.(services.ConflictError); ok {
			c.JSON(http.StatusConflict, gin.H{
				"error": "User with this email already exists",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update user",
			"details": err.Error(),
		})
		return
	}

	// Remove sensitive fields from response
	updatedUser.PasswordHash = nil
	updatedUser.APIKeyHash = nil

	c.JSON(http.StatusOK, gin.H{
		"data":    updatedUser,
		"message": "User updated successfully",
	})
}

// DeleteUser deletes a user
func (h *Handler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	err = h.userService.Delete(c.Request.Context(), id)
	if err != nil {
		if _, ok := err.(services.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete user",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}

// ChangePassword changes user password
func (h *Handler) ChangePassword(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	err = h.userService.ChangePassword(c.Request.Context(), id, req.OldPassword, req.NewPassword)
	if err != nil {
		if _, ok := err.(services.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to change password",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password changed successfully",
	})
}

// GetUsersByCompany gets users by company ID
func (h *Handler) GetUsersByCompany(c *gin.Context) {
	companyIDStr := c.Param("companyId")
	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid company ID",
		})
		return
	}

	users, err := h.userService.GetByCompanyID(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get users by company",
			"details": err.Error(),
		})
		return
	}

	// Remove sensitive fields from response
	for _, user := range users {
		user.PasswordHash = nil
	}

	c.JSON(http.StatusOK, gin.H{
		"data": users,
	})
}

// GetUsersByRole gets users by role
func (h *Handler) GetUsersByRole(c *gin.Context) {
	role := c.Param("role")

	users, err := h.userService.GetByRole(c.Request.Context(), role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get users by role",
			"details": err.Error(),
		})
		return
	}

	// Remove sensitive fields from response
	for _, user := range users {
		user.PasswordHash = nil
	}

	c.JSON(http.StatusOK, gin.H{
		"data": users,
	})
}
