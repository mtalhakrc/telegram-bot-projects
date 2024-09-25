package service

import (
	"context"
	"github.com/haytek-uni-bot-yeniden/pkg/model"
	"github.com/uptrace/bun"
)

type IBaseService[T model.IModel] interface {
	Create(ctx context.Context, t *T) error

	GetByID(ctx context.Context, id int64, relations ...string) (T, error)

	Update(ctx context.Context, t T) error

	Delete(ctx context.Context, id int64) error
	DeleteByUserID(ctx context.Context, id int64) error
}

type BaseService[T model.IModel] struct {
	db *bun.DB
}

func (s BaseService[T]) Create(ctx context.Context, t *T) error {
	_, err := s.db.NewInsert().Model(t).Exec(ctx)
	return err
}

func (s BaseService[T]) GetByID(ctx context.Context, id int64, relations ...string) (T, error) {
	var t T
	sq := s.db.NewSelect().Model(&t).Where("id = ?", id)

	for _, relation := range relations {
		sq = sq.Relation(relation)
	}

	err := sq.Scan(ctx)
	return t, err
}

func (s BaseService[T]) Update(ctx context.Context, t T) error {
	_, err := s.db.NewUpdate().
		Model(&t).
		WherePK().
		Exec(ctx)
	return err
}

func (s BaseService[T]) Delete(ctx context.Context, id int64) error {
	var t T
	_, err := s.db.NewDelete().Model(&t).Where("id = ?", id).Exec(ctx)
	return err
}
func (s BaseService[T]) DeleteByUserID(ctx context.Context, id int64) error {
	var t T
	_, err := s.db.NewDelete().Model(&t).Where("user_id = ?", id).Exec(ctx)
	return err
}

func NewBaseService[T model.IModel](db *bun.DB) IBaseService[T] {
	return BaseService[T]{
		db: db,
	}
}
