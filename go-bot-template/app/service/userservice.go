package service

import (
	"context"
	"github.com/go-bot-template/common/model"
	"github.com/uptrace/bun"
)

type IUserService interface {
	GetByUserID(ctx context.Context, id int64) (model.User, error)

	GetAllUserIDs(ctx context.Context) ([]int64, error)
}
type UserService struct {
	Db *bun.DB
}

func (s UserService) GetAllUserIDs(ctx context.Context) ([]int64, error) {
	var ids []int64
	err := s.Db.NewSelect().Model(&model.User{}).Column("user_id").Scan(ctx, &ids)
	return ids, err
}
func (s UserService) GetByUserID(ctx context.Context, id int64) (model.User, error) {
	u := model.User{}
	err := s.Db.NewSelect().Model(&u).Where("user_id = ?", id).Scan(ctx)
	return u, err
}

func NewUserService(db *bun.DB) IUserService {
	return UserService{
		Db: db,
	}
}
