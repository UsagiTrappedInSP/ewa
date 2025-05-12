package service

import (
	"context"
	"ewa/internal/model"
	"ewa/internal/repository"
	"fmt"
	"time"
)

type UserService struct {
	userRepo     *repository.UserRepository
	redisService RedisService
}

func NewUserService(userRepo *repository.UserRepository, redisService RedisService) *UserService {
	return &UserService{
		userRepo:     userRepo,
		redisService: redisService,
	}
}

func (s *UserService) Create(user *model.User) error {
	if err := s.userRepo.Create(user); err != nil {
		return err
	}

	// 缓存用户信息
	if err := s.redisService.SetUserInfo(context.Background(), fmt.Sprintf("%d", user.ID), user); err != nil {
		// 记录错误但不影响主流程
		// TODO: 添加日志
	}

	return nil
}

func (s *UserService) GetByID(id string) (*model.User, error) {
	// 尝试从Redis获取用户信息
	if user, err := s.redisService.GetUserInfo(context.Background(), id); err == nil {
		return user, nil
	}

	// 从数据库获取用户信息
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 缓存用户信息
	if err := s.redisService.SetUserInfo(context.Background(), id, user); err != nil {
		// 记录错误但不影响主流程
		// TODO: 添加日志
	}

	return user, nil
}

func (s *UserService) GetByEmail(email string) (*model.User, error) {
	// 从数据库获取用户信息
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	// 缓存用户信息
	if err := s.redisService.SetUserInfo(context.Background(), fmt.Sprintf("%d", user.ID), user); err != nil {
		// 记录错误但不影响主流程
		// TODO: 添加日志
	}

	return user, nil
}

func (s *UserService) Update(user *model.User) error {
	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	// 更新缓存
	if err := s.redisService.SetUserInfo(context.Background(), fmt.Sprintf("%d", user.ID), user); err != nil {
		// 记录错误但不影响主流程
		// TODO: 添加日志
	}

	return nil
}

func (s *UserService) Delete(id string) error {
	if err := s.userRepo.Delete(id); err != nil {
		return err
	}

	// 删除缓存
	if err := s.redisService.DeleteUserInfo(context.Background(), id); err != nil {
		// 记录错误但不影响主流程
		// TODO: 添加日志
	}

	return nil
}

func (s *UserService) CheckRateLimit(userID string) bool {
	return s.redisService.IsRateLimited(context.Background(), userID)
} 