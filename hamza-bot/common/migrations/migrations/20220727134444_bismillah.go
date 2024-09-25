package migrations

import (
	"context"
	migration "github.com/haytek-uni-bot-yeniden/common/migrations/migrations/20220727134444_bismillah"
	"github.com/uptrace/bun"
)

func init() {
	m := []interface{}{
		&migration.User{},
		&migration.GunlukRapor{},
	}
	up := func(ctx context.Context, db *bun.DB) error {
		// manytomany için önce register yapmak gerekiyormuş.
		// manytomany için önce register yapmak gerekiyormuş.
		//db.RegisterModel((*migration.Kitap)(nil))
		//db.RegisterModel((*migration.Section)(nil))
		return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
			for _, i := range m {
				if _, err := tx.NewCreateTable().Model(i).IfNotExists().WithForeignKeys().Exec(ctx); err != nil {
					return err
				}
			}
			return nil
		})
	}

	down := func(ctx context.Context, db *bun.DB) error {
		return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
			// manytomany için önce register yapmak gerekiyormuş.
			//db.RegisterModel((*migration.Kitap)(nil))
			//db.RegisterModel((*migration.Section)(nil))
			for _, i := range m {
				if _, err := tx.NewDropTable().Model(i).IfExists().Cascade().Exec(ctx); err != nil {
					return err
				}
			}
			return nil
		})
	}

	if err := Migrations.Register(up, down); err != nil {
		panic(err)
	}
}
