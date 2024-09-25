package fixture

import (
	"github.com/haytek-uni-bot-yeniden/common/model"
	"github.com/uptrace/bun"
)

func Load(db *bun.DB) error {
	db.RegisterModel((*model.User)(nil))

	//fixture := dbfixture.New(db, dbfixture.WithTruncateTables())
	//err := fixture.Load(context.Background(), os.DirFS("common/fixture"), "fixture.yml")

	//return err
	return nil
}
