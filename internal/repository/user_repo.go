package repository

import (
	"test/internal/domain"
	"test/pkg/db"
	"test/pkg/logger"
)

type UserRepository struct {
	wrapper *db.MySqlWrapper
}

func NewUserRepository(db *db.MySqlWrapper) (*UserRepository, error) {
	return &UserRepository{
		wrapper: db,
	}, nil
}

// CRUD操作
func (u *UserRepository) GetAll() ([]*domain.User, error) {
	var users []*domain.User
	res := u.wrapper.DB.Find(&users)
	if res.Error != nil {
		logger.Error("查询所有用户失败!")
		return nil, res.Error
	}
	logger.Info("查询所有用户成功...")
	return users, nil
}
