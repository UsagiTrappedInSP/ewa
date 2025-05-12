package service

import (
	"context"
	"ewa/internal/model"
	"ewa/internal/repository"
	"fmt"
	"strconv"
	"time"
)

type ModelService struct {
	modelRepo      *repository.ModelRepository
	permissionRepo *repository.PermissionRepository
	ossService     OSSService
	redisService   RedisService
}

func NewModelService(modelRepo *repository.ModelRepository, permissionRepo *repository.PermissionRepository, ossService OSSService, redisService RedisService) *ModelService {
	return &ModelService{
		modelRepo:      modelRepo,
		permissionRepo: permissionRepo,
		ossService:     ossService,
		redisService:   redisService,
	}
}

func (s *ModelService) Create(model *model.Model) error {
	// 创建模型
	if err := s.modelRepo.Create(model); err != nil {
		return err
	}

	// 为创建者添加超级权限
	permission := &model.UserModelPermission{
		UserID:          model.CreatedBy,
		ModelID:         model.ID,
		PermissionLevel: 4, // 超级权限
		IsOwner:         1, // 所有者
		GrantedBy:       model.CreatedBy,
		Status:          1,
	}

	if err := s.permissionRepo.Create(permission); err != nil {
		return err
	}

	// 如果是公共模型（user_id=1），添加到Redis缓存
	if model.CreatedBy == 1 {
		if err := s.redisService.AddPublicModel(context.Background(), model); err != nil {
			// 记录错误但不影响主流程
			// TODO: 添加日志
		}
	}

	return nil
}

func (s *ModelService) GetByID(id string) (*model.Model, error) {
	// 尝试从Redis获取公共模型
	if models, err := s.redisService.GetPublicModels(context.Background()); err == nil {
		for _, m := range models {
			if fmt.Sprintf("%d", m.ID) == id {
				return m, nil
			}
		}
	}

	return s.modelRepo.GetByID(id)
}

func (s *ModelService) Update(model *model.Model) error {
	if err := s.modelRepo.Update(model); err != nil {
		return err
	}

	// 如果是公共模型，更新Redis缓存
	if model.CreatedBy == 1 {
		if err := s.redisService.RemovePublicModel(context.Background(), fmt.Sprintf("%d", model.ID)); err != nil {
			// 记录错误但不影响主流程
			// TODO: 添加日志
		}
		if err := s.redisService.AddPublicModel(context.Background(), model); err != nil {
			// 记录错误但不影响主流程
			// TODO: 添加日志
		}
	}

	return nil
}

func (s *ModelService) Delete(id string) error {
	// 获取模型信息
	model, err := s.modelRepo.GetByID(id)
	if err != nil {
		return err
	}

	if err := s.modelRepo.Delete(id); err != nil {
		return err
	}

	// 如果是公共模型，从Redis缓存中删除
	if model.CreatedBy == 1 {
		if err := s.redisService.RemovePublicModel(context.Background(), id); err != nil {
			// 记录错误但不影响主流程
			// TODO: 添加日志
		}
	}

	return nil
}

func (s *ModelService) List(userID uint64, page, pageSize, category, keyword string) ([]*model.Model, int64, error) {
	pageNum, _ := strconv.Atoi(page)
	pageSizeNum, _ := strconv.Atoi(pageSize)

	// 获取用户有权限的模型ID列表
	modelIDs, err := s.permissionRepo.GetModelIDsByUserID(userID)
	if err != nil {
		return nil, 0, err
	}

	// 如果是第一页且没有筛选条件，尝试从Redis获取公共模型
	if pageNum == 1 && category == "" && keyword == "" {
		if models, err := s.redisService.GetPublicModels(context.Background()); err == nil {
			// 过滤出用户有权限的模型
			var filteredModels []*model.Model
			for _, m := range models {
				for _, id := range modelIDs {
					if m.ID == id {
						filteredModels = append(filteredModels, m)
						break
					}
				}
			}
			if len(filteredModels) > 0 {
				return filteredModels, int64(len(filteredModels)), nil
			}
		}
	}

	return s.modelRepo.List(modelIDs, pageNum, pageSizeNum, category, keyword)
}

func (s *ModelService) CheckPermission(userID uint64, modelID string, requiredLevel int8) (bool, error) {
	// 尝试从Redis获取权限信息
	if permissions, err := s.redisService.GetModelPermissions(context.Background(), modelID); err == nil {
		for _, p := range permissions {
			if p.UserID == userID && p.Status == 1 {
				// 检查权限是否过期
				if p.ExpireTime != nil && p.ExpireTime.Before(time.Now()) {
					return false, nil
				}
				return p.PermissionLevel >= requiredLevel, nil
			}
		}
	}

	// 从数据库获取权限信息
	permission, err := s.permissionRepo.GetByUserIDAndModelID(userID, modelID)
	if err != nil {
		return false, err
	}

	// 检查权限是否过期
	if permission.ExpireTime != nil && permission.ExpireTime.Before(time.Now()) {
		return false, nil
	}

	// 更新Redis缓存
	if err := s.redisService.SetModelPermissions(context.Background(), modelID, []*model.UserModelPermission{permission}); err != nil {
		// 记录错误但不影响主流程
		// TODO: 添加日志
	}

	return permission.PermissionLevel >= requiredLevel, nil
}

// GetModelFile 获取模型文件信息
func (s *ModelService) GetModelFile(ctx context.Context, modelID string) (*model.ModelFileResponse, error) {
	model, err := s.modelRepo.GetByID(modelID)
	if err != nil {
		return nil, err
	}

	// 获取文件的临时访问URL
	fileURL, err := s.ossService.GetFileURL(ctx, model.FileURL)
	if err != nil {
		return nil, err
	}

	return &model.ModelFileResponse{
		FileURL:    fileURL,
		FileName:   model.Name + "." + model.FileType,
		FileSize:   model.FileSize,
		FileType:   model.FileType,
		ExpireTime: time.Now().Add(time.Hour).Format(time.RFC3339),
	}, nil
}

// GetUploadURL 获取文件上传URL
func (s *ModelService) GetUploadURL(ctx context.Context, fileName, fileType string) (*model.UploadURLResponse, error) {
	objectKey := "models/" + time.Now().Format("2006/01/02") + "/" + fileName
	uploadURL, fileURL, err := s.ossService.GetUploadURL(ctx, objectKey, fileType)
	if err != nil {
		return nil, err
	}

	return &model.UploadURLResponse{
		UploadURL:  uploadURL,
		FileURL:    fileURL,
		ExpireTime: time.Now().Add(time.Hour).Format(time.RFC3339),
	}, nil
} 