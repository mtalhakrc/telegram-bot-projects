package service

import (
	"context"
	"github.com/haytek-uni-bot-yeniden/common/model"
	"github.com/uptrace/bun"
)

type ISessionService interface {
	GetByUserID(ctx context.Context, id int64) (model.Session, error)
}
type SessionService struct {
	IBaseService[model.Session]
	Db *bun.DB
}

func (s SessionService) GetByUserID(ctx context.Context, id int64) (model.Session, error) {
	u := model.Session{}
	err := s.Db.NewSelect().Model(&u).Where("user_id = ?", id).Scan(ctx)
	return u, err
}

func NewSessionService(db *bun.DB) ISessionService {
	return SessionService{
		IBaseService: NewBaseService[model.Session](db),
		Db:           db,
	}
}
