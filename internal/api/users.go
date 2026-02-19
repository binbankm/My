package api

import (
	"net/http"

	"github.com/binbankm/My/internal/model"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// ListUsers lists all users
func ListUsers(c *gin.Context) {
	var users []model.User
	
	if err := model.DB.Preload("Role").Preload("Role.Permissions").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetUser gets a specific user
func GetUser(c *gin.Context) {
	id := c.Param("id")

	var user model.User
	if err := model.DB.Preload("Role").Preload("Role.Permissions").First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// CreateUser creates a new user
func CreateUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Email    string `json:"email"`
		IsAdmin  bool   `json:"isAdmin"`
		RoleID   *uint  `json:"roleId"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if username already exists
	var existing model.User
	if err := model.DB.Where("username = ?", req.Username).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
		IsAdmin:  req.IsAdmin,
		RoleID:   req.RoleID,
		Enabled:  true,
	}

	if err := model.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
		return
	}

	// Load role information
	model.DB.Preload("Role").First(&user, user.ID)

	c.JSON(http.StatusOK, user)
}

// UpdateUser updates a user
func UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		IsAdmin  bool   `json:"isAdmin"`
		RoleID   *uint  `json:"roleId"`
		Enabled  bool   `json:"enabled"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user model.User
	if err := model.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Update fields
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		user.Password = string(hashedPassword)
	}
	user.IsAdmin = req.IsAdmin
	user.RoleID = req.RoleID
	user.Enabled = req.Enabled

	if err := model.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user: " + err.Error()})
		return
	}

	// Load role information
	model.DB.Preload("Role").First(&user, user.ID)

	c.JSON(http.StatusOK, user)
}

// DeleteUser deletes a user
func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	var user model.User
	if err := model.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Don't allow deleting the last admin
	if user.IsAdmin {
		var adminCount int64
		model.DB.Model(&model.User{}).Where("is_admin = ?", true).Count(&adminCount)
		if adminCount <= 1 {
			c.JSON(http.StatusForbidden, gin.H{"error": "Cannot delete the last admin user"})
			return
		}
	}

	if err := model.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// ListRoles lists all roles
func ListRoles(c *gin.Context) {
	var roles []model.Role
	
	if err := model.DB.Preload("Permissions").Find(&roles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch roles: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, roles)
}

// GetRole gets a specific role
func GetRole(c *gin.Context) {
	id := c.Param("id")

	var role model.Role
	if err := model.DB.Preload("Permissions").First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	c.JSON(http.StatusOK, role)
}

// CreateRole creates a new role
func CreateRole(c *gin.Context) {
	var req struct {
		Name          string `json:"name" binding:"required"`
		Description   string `json:"description"`
		PermissionIDs []uint `json:"permissionIds"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get permissions
	var permissions []model.Permission
	if len(req.PermissionIDs) > 0 {
		if err := model.DB.Find(&permissions, req.PermissionIDs).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid permission IDs"})
			return
		}
	}

	role := model.Role{
		Name:        req.Name,
		Description: req.Description,
		Permissions: permissions,
	}

	if err := model.DB.Create(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create role: " + err.Error()})
		return
	}

	// Load permissions
	model.DB.Preload("Permissions").First(&role, role.ID)

	c.JSON(http.StatusOK, role)
}

// UpdateRole updates a role
func UpdateRole(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name          string `json:"name"`
		Description   string `json:"description"`
		PermissionIDs []uint `json:"permissionIds"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var role model.Role
	if err := model.DB.First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	// Update fields
	if req.Name != "" {
		role.Name = req.Name
	}
	if req.Description != "" {
		role.Description = req.Description
	}

	// Update permissions
	if len(req.PermissionIDs) > 0 {
		var permissions []model.Permission
		if err := model.DB.Find(&permissions, req.PermissionIDs).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid permission IDs"})
			return
		}
		
		// Clear existing permissions and set new ones
		if err := model.DB.Model(&role).Association("Permissions").Replace(&permissions); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update permissions"})
			return
		}
	}

	if err := model.DB.Save(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update role: " + err.Error()})
		return
	}

	// Load permissions
	model.DB.Preload("Permissions").First(&role, role.ID)

	c.JSON(http.StatusOK, role)
}

// DeleteRole deletes a role
func DeleteRole(c *gin.Context) {
	id := c.Param("id")

	var role model.Role
	if err := model.DB.First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	// Check if any users are using this role
	var userCount int64
	model.DB.Model(&model.User{}).Where("role_id = ?", id).Count(&userCount)
	if userCount > 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot delete role: it is assigned to users"})
		return
	}

	if err := model.DB.Delete(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete role: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
}

// ListPermissions lists all permissions
func ListPermissions(c *gin.Context) {
	var permissions []model.Permission
	
	if err := model.DB.Find(&permissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch permissions: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, permissions)
}
