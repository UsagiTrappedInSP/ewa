package handler

import (
	"net/http"
	"ewa/internal/model"
	"ewa/internal/service"
	"github.com/gin-gonic/gin"
)

type PermissionHandler struct {
	permissionService *service.PermissionService
}

func NewPermissionHandler(permissionService *service.PermissionService) *PermissionHandler {
	return &PermissionHandler{
		permissionService: permissionService,
	}
}

// Grant 授予权限
func (h *PermissionHandler) Grant(c *gin.Context) {
	var req model.PermissionGrantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查授权人是否有权限
	grantedBy := c.GetUint64("user_id")
	hasPermission, err := h.permissionService.CheckPermission(grantedBy, req.ModelID, 3) // 需要管理权限
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	permission := &model.UserModelPermission{
		UserID:          req.UserID,
		ModelID:         req.ModelID,
		PermissionLevel: req.PermissionLevel,
		IsOwner:         0,
		ExpireTime:      req.ExpireTime,
		GrantedBy:       grantedBy,
		GrantReason:     req.GrantReason,
		Status:          1,
	}

	if req.IsOwner != nil && *req.IsOwner == 1 {
		permission.IsOwner = 1
	}

	if err := h.permissionService.Create(permission); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Permission granted successfully",
		"permission": permission,
	})
}

// Update 更新权限
func (h *PermissionHandler) Update(c *gin.Context) {
	permissionID := c.Param("id")
	var req model.PermissionUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查授权人是否有权限
	grantedBy := c.GetUint64("user_id")
	permission, err := h.permissionService.GetByID(permissionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Permission not found"})
		return
	}

	hasPermission, err := h.permissionService.CheckPermission(grantedBy, permission.ModelID, 3) // 需要管理权限
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	// 更新字段
	if req.PermissionLevel != 0 {
		permission.PermissionLevel = req.PermissionLevel
	}
	if req.IsOwner != nil {
		permission.IsOwner = *req.IsOwner
	}
	if req.ExpireTime != nil {
		permission.ExpireTime = req.ExpireTime
	}
	if req.Status != nil {
		permission.Status = *req.Status
	}

	if err := h.permissionService.Update(permission); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Permission updated successfully",
		"permission": permission,
	})
}

// Revoke 撤销权限
func (h *PermissionHandler) Revoke(c *gin.Context) {
	permissionID := c.Param("id")

	// 检查授权人是否有权限
	grantedBy := c.GetUint64("user_id")
	permission, err := h.permissionService.GetByID(permissionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Permission not found"})
		return
	}

	hasPermission, err := h.permissionService.CheckPermission(grantedBy, permission.ModelID, 3) // 需要管理权限
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	if err := h.permissionService.Delete(permissionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permission revoked successfully"})
}

// List 获取权限列表
func (h *PermissionHandler) List(c *gin.Context) {
	modelID := c.Param("model_id")
	userID := c.GetUint64("user_id")

	// 检查权限
	hasPermission, err := h.permissionService.CheckPermission(userID, modelID, 3) // 需要管理权限
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	permissions, err := h.permissionService.ListByModelID(modelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"permissions": permissions})
} 