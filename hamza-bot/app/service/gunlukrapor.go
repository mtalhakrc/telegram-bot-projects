package service

import (
	"context"
	"github.com/haytek-uni-bot-yeniden/common/model"
	"github.com/haytek-uni-bot-yeniden/pkg/service"
	"github.com/uptrace/bun"
)

type GunlukRaporService struct {
	db *bun.DB
	service.IBaseService[model.GunlukRapor]
}

type IGunlukRaporService interface {
	Kaydet(m model.GunlukRapor) error
	GetLastWeekRecords() ([]model.GunlukRapor, error)
	GetSpecificDayRecord(date string) (model.GunlukRapor, error)
}

func NewGunlukRaporService(db *bun.DB) IGunlukRaporService {
	return GunlukRaporService{
		IBaseService: service.NewBaseService[model.GunlukRapor](db),
		db:           db,
	}
}

func (s GunlukRaporService) Kaydet(m model.GunlukRapor) error {
	return s.Create(context.Background(), &m)
}
func (s GunlukRaporService) GetLastWeekRecords() ([]model.GunlukRapor, error) {
	var m []model.GunlukRapor
	err := s.db.NewSelect().Model(&m).Where("created_at  >= datetime('now', '-6 days')").Scan(context.Background())
	return m, err
}
func (s GunlukRaporService) GetSpecificDayRecord(date string) (model.GunlukRapor, error) {
	var m model.GunlukRapor
	err := s.db.NewSelect().Model(&m).Where("created_at >= date(?)", date).Scan(context.Background())
	return m, err
}
