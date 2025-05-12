package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"ewa/internal/model"
	"github.com/redis/go-redis/v9"
)

// RedisService 定义Redis服务接口
type RedisService interface {
	// 公共模型相关
	GetPublicModels(ctx context.Context) ([]*model.Model, error)
	SetPublicModels(ctx context.Context, models []*model.Model) error
	AddPublicModel(ctx context.Context, model *model.Model) error
	RemovePublicModel(ctx context.Context, modelID string) error

	// 用户信息相关
	GetUserInfo(ctx context.Context, userID uint64) (*model.User, error)
	SetUserInfo(ctx context.Context, user *model.User) error
	DeleteUserInfo(ctx context.Context, userID uint64) error

	// 模型权限相关
	GetModelPermissions(ctx context.Context, modelID string) ([]*model.UserModelPermission, error)
	SetModelPermissions(ctx context.Context, modelID string, permissions []*model.UserModelPermission) error
	DeleteModelPermissions(ctx context.Context, modelID string) error

	// 限流相关
	IsRateLimited(ctx context.Context, key string, limit int) (bool, error)
}

// RedisServiceImpl Redis服务实现
type RedisServiceImpl struct {
	client *redis.Client
	config *RedisConfig
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host            string
	Port            int
	Password        string
	DB              int
	PoolSize        int
	MinIdleConns    int
	CacheExpiration map[string]time.Duration
}

// NewRedisService 创建Redis服务实例
func NewRedisService(config *RedisConfig) (RedisService, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password:     config.Password,
		DB:           config.DB,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
	})

	// 测试连接
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %v", err)
	}

	return &RedisServiceImpl{
		client: client,
		config: config,
	}, nil
}

// GetPublicModels 获取公共模型列表
func (s *RedisServiceImpl) GetPublicModels(ctx context.Context) ([]*model.Model, error) {
	data, err := s.client.Get(ctx, "public_models").Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var models []*model.Model
	if err := json.Unmarshal(data, &models); err != nil {
		return nil, err
	}
	return models, nil
}

// SetPublicModels 设置公共模型列表
func (s *RedisServiceImpl) SetPublicModels(ctx context.Context, models []*model.Model) error {
	data, err := json.Marshal(models)
	if err != nil {
		return err
	}

	return s.client.Set(ctx, "public_models", data, s.config.CacheExpiration["public_models"]).Err()
}

// AddPublicModel 添加公共模型
func (s *RedisServiceImpl) AddPublicModel(ctx context.Context, model *model.Model) error {
	models, err := s.GetPublicModels(ctx)
	if err != nil {
		return err
	}

	models = append(models, model)
	return s.SetPublicModels(ctx, models)
}

// RemovePublicModel 移除公共模型
func (s *RedisServiceImpl) RemovePublicModel(ctx context.Context, modelID string) error {
	models, err := s.GetPublicModels(ctx)
	if err != nil {
		return err
	}

	for i, m := range models {
		if fmt.Sprintf("%d", m.ID) == modelID {
			models = append(models[:i], models[i+1:]...)
			break
		}
	}

	return s.SetPublicModels(ctx, models)
}

// GetUserInfo 获取用户信息
func (s *RedisServiceImpl) GetUserInfo(ctx context.Context, userID uint64) (*model.User, error) {
	data, err := s.client.Get(ctx, fmt.Sprintf("user:%d", userID)).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var user model.User
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// SetUserInfo 设置用户信息
func (s *RedisServiceImpl) SetUserInfo(ctx context.Context, user *model.User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return s.client.Set(ctx, fmt.Sprintf("user:%d", user.ID), data, s.config.CacheExpiration["user_info"]).Err()
}

// DeleteUserInfo 删除用户信息
func (s *RedisServiceImpl) DeleteUserInfo(ctx context.Context, userID uint64) error {
	return s.client.Del(ctx, fmt.Sprintf("user:%d", userID)).Err()
}

// GetModelPermissions 获取模型权限列表
func (s *RedisServiceImpl) GetModelPermissions(ctx context.Context, modelID string) ([]*model.UserModelPermission, error) {
	data, err := s.client.Get(ctx, fmt.Sprintf("model_permissions:%s", modelID)).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var permissions []*model.UserModelPermission
	if err := json.Unmarshal(data, &permissions); err != nil {
		return nil, err
	}
	return permissions, nil
}

// SetModelPermissions 设置模型权限列表
func (s *RedisServiceImpl) SetModelPermissions(ctx context.Context, modelID string, permissions []*model.UserModelPermission) error {
	data, err := json.Marshal(permissions)
	if err != nil {
		return err
	}

	return s.client.Set(ctx, fmt.Sprintf("model_permissions:%s", modelID), data, s.config.CacheExpiration["model_permissions"]).Err()
}

// DeleteModelPermissions 删除模型权限列表
func (s *RedisServiceImpl) DeleteModelPermissions(ctx context.Context, modelID string) error {
	return s.client.Del(ctx, fmt.Sprintf("model_permissions:%s", modelID)).Err()
}

// IsRateLimited 检查是否超过速率限制
func (s *RedisServiceImpl) IsRateLimited(ctx context.Context, key string, limit int) (bool, error) {
	count, err := s.client.Incr(ctx, fmt.Sprintf("rate_limit:%s", key)).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		s.client.Expire(ctx, fmt.Sprintf("rate_limit:%s", key), s.config.CacheExpiration["rate_limit"])
	}

	return count > int64(limit), nil
} 