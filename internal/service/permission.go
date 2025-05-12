package service

import (
	"context"
	"ewa/internal/model"
	"ewa/internal/repository"
	"fmt"
	"time"
)

type PermissionService struct {
	permissionRepo *repository.PermissionRepository
	redisService   RedisService
}

func NewPermissionService(permissionRepo *repository.PermissionRepository, redisService RedisService) *PermissionService {
	return &PermissionService{
		permissionRepo: permissionRepo,
		redisService:   redisService,
	}
}

func (s *PermissionService) Create(permission *model.UserModelPermission) error {
	if err := s.permissionRepo.Create(permission); err != nil {
		return err
	}

	// 更新Redis缓存
	if err := s.redisService.SetModelPermissions(context.Background(), fmt.Sprintf("%d", permission.ModelID), []*model.UserModelPermission{permission}); err != nil {
		// 记录错误但不影响主流程
		// TODO: 添加日志
	}

	return nil
}

func (s *PermissionService) GetByID(id string) (*model.UserModelPermission, error) {
	return s.permissionRepo.GetByID(id)
}

func (s *PermissionService) Update(permission *model.UserModelPermission) error {
	if err := s.permissionRepo.Update(permission); err != nil {
		return err
	}

	// 更新Redis缓存
	if err := s.redisService.SetModelPermissions(context.Background(), fmt.Sprintf("%d", permission.ModelID), []*model.UserModelPermission{permission}); err != nil {
		// 记录错误但不影响主流程
		// TODO: 添加日志
	}

	return nil
}

func (s *PermissionService) Delete(id string) error {
	// 获取权限信息
	permission, err := s.permissionRepo.GetByID(id)
	if err != nil {
		return err
	}

	if err := s.permissionRepo.Delete(id); err != nil {
		return err
	}

	// 更新Redis缓存
	if err := s.redisService.DeleteModelPermissions(context.Background(), fmt.Sprintf("%d", permission.ModelID)); err != nil {
		// 记录错误但不影响主流程
		// TODO: 添加日志
	}

	return nil
}

func (s *PermissionService) ListByModelID(modelID string) ([]*model.UserModelPermission, error) {
	// 尝试从Redis获取权限信息
	if permissions, err := s.redisService.GetModelPermissions(context.Background(), modelID); err == nil {
		return permissions, nil
	}

	// 从数据库获取权限信息
	permissions, err := s.permissionRepo.ListByModelID(modelID)
	if err != nil {
		return nil, err
	}

	// 更新Redis缓存
	if err := s.redisService.SetModelPermissions(context.Background(), modelID, permissions); err != nil {
		// 记录错误但不影响主流程
		// TODO: 添加日志
	}

	return permissions, nil
}

func (s *PermissionService) GetByUserIDAndModelID(userID uint64, modelID string) (*model.UserModelPermission, error) {
	// 尝试从Redis获取权限信息
	if permissions, err := s.redisService.GetModelPermissions(context.Background(), modelID); err == nil {
		for _, p := range permissions {
			if p.UserID == userID && p.Status == 1 {
				// 检查权限是否过期
				if p.ExpireTime != nil && p.ExpireTime.Before(time.Now()) {
					return nil, nil
				}
				return p, nil
			}
		}
	}

	// 从数据库获取权限信息
	permission, err := s.permissionRepo.GetByUserIDAndModelID(userID, modelID)
	if err != nil {
		return nil, err
	}

	// 更新Redis缓存
	if err := s.redisService.SetModelPermissions(context.Background(), modelID, []*model.UserModelPermission{permission}); err != nil {
		// 记录错误但不影响主流程
		// TODO: 添加日志
	}

	return permission, nil
}

func (s *PermissionService) CheckPermission(userID uint64, modelID string, requiredLevel int8) (bool, error) {
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