package service

import (
	"context"
	"go.uber.org/zap"
	"test/internal/domain"
	"test/internal/repository"
	"test/pkg/logger"
)

type UserService struct {
	UserRepository *repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository) (*UserService, error) {

	return &UserService{
		UserRepository: userRepository,
	}, nil
}
func (s *UserService) GetUserByID(ctx context.Context, userId int64) (*domain.User, error) {
	logger.Info("查询用户", zap.Int64("user_id", userId))
	//数据库查询业务
	return &domain.User{
		UserID:   "",
		Password: "",
		UserName: "",
		Gender:   0,
		Account:  "",
		Email:    "",
		Role:     0,
	}, nil
}
func (s *UserService) GetALLUsers() ([]*domain.User, error) {
	all, err := s.UserRepository.GetAll()
	if err != nil {
		return nil, err
	}
	return all, nil
}
