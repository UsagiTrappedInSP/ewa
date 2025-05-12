package handler

import (
	"net/http"
	"ewa/internal/model"
	"ewa/internal/service"
	"github.com/gin-gonic/gin"
)

type ModelHandler struct {
	modelService *service.ModelService
}

func NewModelHandler(modelService *service.ModelService) *ModelHandler {
	return &ModelHandler{
		modelService: modelService,
	}
}

// Create 创建模型
func (h *ModelHandler) Create(c *gin.Context) {
	var req model.ModelCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint64("user_id")
	model := &model.Model{
		Name:        req.Name,
		Version:     req.Version,
		FileURL:     req.FileURL,
		FileSize:    req.FileSize,
		FileType:    req.FileType,
		Description: req.Description,
		Category:    req.Category,
		Tags:        req.Tags,
		Status:      1,
		CreatedBy:   userID,
	}

	if err := h.modelService.Create(model); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Model created successfully",
		"model":   model,
	})
}

// Update 更新模型
func (h *ModelHandler) Update(c *gin.Context) {
	modelID := c.Param("id")
	var req model.ModelUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查权限
	userID := c.GetUint64("user_id")
	hasPermission, err := h.modelService.CheckPermission(userID, modelID, 2) // 需要编辑权限
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	model, err := h.modelService.GetByID(modelID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Model not found"})
		return
	}

	// 更新字段
	if req.Name != "" {
		model.Name = req.Name
	}
	if req.Version != "" {
		model.Version = req.Version
	}
	if req.Description != "" {
		model.Description = req.Description
	}
	if req.Category != "" {
		model.Category = req.Category
	}
	if req.Tags != "" {
		model.Tags = req.Tags
	}
	if req.Status != nil {
		model.Status = *req.Status
	}

	if err := h.modelService.Update(model); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Model updated successfully",
		"model":   model,
	})
}

// Delete 删除模型
func (h *ModelHandler) Delete(c *gin.Context) {
	modelID := c.Param("id")
	userID := c.GetUint64("user_id")

	// 检查权限
	hasPermission, err := h.modelService.CheckPermission(userID, modelID, 4) // 需要超级权限
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	if err := h.modelService.Delete(modelID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Model deleted successfully"})
}

// Get 获取模型详情
func (h *ModelHandler) Get(c *gin.Context) {
	modelID := c.Param("id")
	userID := c.GetUint64("user_id")

	// 检查权限
	hasPermission, err := h.modelService.CheckPermission(userID, modelID, 1) // 需要查看权限
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	model, err := h.modelService.GetByID(modelID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Model not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"model": model})
}

// List 获取模型列表
func (h *ModelHandler) List(c *gin.Context) {
	userID := c.GetUint64("user_id")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "10")
	category := c.Query("category")
	keyword := c.Query("keyword")

	models, total, err := h.modelService.List(userID, page, pageSize, category, keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"models": models,
		"total":  total,
	})
}

// GetModelFile 获取模型文件
func (h *ModelHandler) GetModelFile(c *gin.Context) {
	modelID := c.Param("id")
	userID := c.GetUint64("user_id")

	// 检查权限
	hasPermission, err := h.modelService.CheckPermission(userID, modelID, 1) // 需要查看权限
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	fileInfo, err := h.modelService.GetModelFile(c.Request.Context(), modelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, fileInfo)
}

// GetUploadURL 获取文件上传URL
func (h *ModelHandler) GetUploadURL(c *gin.Context) {
	var req model.UploadURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uploadInfo, err := h.modelService.GetUploadURL(c.Request.Context(), req.FileName, req.FileType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, uploadInfo)
} 