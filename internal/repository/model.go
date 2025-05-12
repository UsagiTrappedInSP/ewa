package repository

import (
	"ewa/internal/model"
	"gorm.io/gorm"
)

type ModelRepository struct {
	db *gorm.DB
}

func NewModelRepository(db *gorm.DB) *ModelRepository {
	return &ModelRepository{
		db: db,
	}
}

func (r *ModelRepository) Create(model *model.Model) error {
	return r.db.Create(model).Error
}

func (r *ModelRepository) GetByID(id string) (*model.Model, error) {
	var model model.Model
	err := r.db.First(&model, id).Error
	if err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *ModelRepository) Update(model *model.Model) error {
	return r.db.Save(model).Error
}

func (r *ModelRepository) Delete(id string) error {
	return r.db.Delete(&model.Model{}, id).Error
}

func (r *ModelRepository) List(modelIDs []uint64, page, pageSize int, category, keyword string) ([]*model.Model, int64, error) {
	var models []*model.Model
	var total int64

	query := r.db.Model(&model.Model{}).Where("id IN ?", modelIDs)

	if category != "" {
		query = query.Where("category = ?", category)
	}
	if keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	return models, total, nil
} 