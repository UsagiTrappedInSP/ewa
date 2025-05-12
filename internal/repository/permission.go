package repository

import (
	"ewa/internal/model"
	"gorm.io/gorm"
)

type PermissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{
		db: db,
	}
}

func (r *PermissionRepository) Create(permission *model.UserModelPermission) error {
	return r.db.Create(permission).Error
}

func (r *PermissionRepository) GetByID(id string) (*model.UserModelPermission, error) {
	var permission model.UserModelPermission
	err := r.db.First(&permission, id).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *PermissionRepository) Update(permission *model.UserModelPermission) error {
	return r.db.Save(permission).Error
}

func (r *PermissionRepository) Delete(id string) error {
	return r.db.Delete(&model.UserModelPermission{}, id).Error
}

func (r *PermissionRepository) ListByModelID(modelID string) ([]*model.UserModelPermission, error) {
	var permissions []*model.UserModelPermission
	err := r.db.Where("model_id = ?", modelID).Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r *PermissionRepository) GetByUserIDAndModelID(userID uint64, modelID string) (*model.UserModelPermission, error) {
	var permission model.UserModelPermission
	err := r.db.Where("user_id = ? AND model_id = ?", userID, modelID).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *PermissionRepository) GetModelIDsByUserID(userID uint64) ([]uint64, error) {
	var modelIDs []uint64
	err := r.db.Model(&model.UserModelPermission{}).
		Where("user_id = ? AND status = 1", userID).
		Pluck("model_id", &modelIDs).Error
	if err != nil {
		return nil, err
	}
	return modelIDs, nil
} 