package service

import (
	"context"
	"github.com/haytek-uni-bot-yeniden/common/model"
	"github.com/uptrace/bun"
	"strings"
)

type IUserService interface {
	GetByUserID(ctx context.Context, id int64) (model.User, error)

	GetAllUserIDs(ctx context.Context) ([]int64, error)
	GetByName(ctx context.Context, name string) (model.User, error)
	GetByUsername(ctx context.Context, username string) (model.User, error)

	GetAdminsNames(ctx context.Context) ([]string, error)
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
func (s UserService) GetByUsername(ctx context.Context, username string) (model.User, error) {
	u := model.User{}
	err := s.Db.NewSelect().Model(&u).Where("username = ?", username).Scan(ctx)
	return u, err
}

func (s UserService) GetByName(ctx context.Context, name string) (model.User, error) {
	u := model.User{}
	err := s.Db.NewSelect().Model(&u).Where("lower(name) = ?", strings.ToLower(name)).Scan(ctx)
	return u, err
}
func (s UserService) GetAdminsNames(ctx context.Context) ([]string, error) {
	var res []string
	err := s.Db.NewSelect().Model(&model.User{}).Where("type = ?", model.UserTypeAdmin).Column("name").Scan(ctx, &res)
	return res, err
}

func NewUserService(db *bun.DB) IUserService {
	return UserService{
		Db: db,
	}
}
