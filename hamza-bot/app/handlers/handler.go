package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/haytek-uni-bot-yeniden/app/service"
	"github.com/haytek-uni-bot-yeniden/common/model"
	"github.com/haytek-uni-bot-yeniden/pkg/app"
	baseservice "github.com/haytek-uni-bot-yeniden/pkg/service"
	"github.com/uptrace/bun"
	"log"
	"strings"
)

type UserHandler struct {
	baseservice.IBaseService[model.User]
	userService   service.IUserService
	sheetsService service.ISheetsService
}

func NewUserHandler(db *bun.DB, s service.ISheetsService) *UserHandler {
	return &UserHandler{
		IBaseService:  baseservice.NewBaseService[model.User](db),
		userService:   service.NewUserService(db),
		sheetsService: s,
	}
}
func (u UserHandler) Kaydol(ctx *app.Ctx, params []string) (string, error) {

	username := ctx.SentFrom().String()
	userid := ctx.SentFrom().ID

	//zaten kaydoldu mu kontrol et
	_, err := u.userService.GetByUserID(context.Background(), userid)
	if err == nil {
		err = errors.New("kaydınız zaten var")
		return "", err
	}

	name := strings.Join(params, " ")

	if name == "" {
		return "", errors.New("komuttan sonra isminizi belirtiniz(ör: /kaydol Talha Karaca)")
	}
	if !u.sheetsService.TestSheetExist(name) {
		return "", errors.New("isminiz çetelede bulunamadı, çetelenizde bu isme ait bir sayfa olduğundan emin olup tekrar deneyiniz")
	}

	user := model.User{
		Name:     name,
		Username: username,
		UserID:   userid,
		Type:     model.UserTypeNormal,
	}

	err = u.Create(context.Background(), &user)
	if err != nil {
		return "", err
	}
	return "kaydınız başarı ile gerçekleştirildi", nil
}

func (u UserHandler) KayitSil(ctx *app.Ctx, params []string) (string, error) {
	err := u.DeleteByUserID(context.Background(), ctx.SentFrom().ID)
	return "kaydınız başarı ile silindi", err
}

func (u UserHandler) MakeAdmin(ctx *app.Ctx, params []string) (string, error) {
	if !u.IsAdmin(ctx.SentFrom().ID) {
		return "", errors.New("bu komut için yetkiniz bulunmamaktadır")
	}
	if len(params) != 1 {
		return "", errors.New("komuttan sonra bir kişiyi belirtin(ör: /makeadmin @mtalhakrc)")
	}
	if params[0][0] != '@' {
		return "", errors.New("kişinin başında @ işareti olmalıdır")
	}

	username := strings.Replace(params[0], "@", "", 1)
	user, err := u.userService.GetByUsername(context.Background(), username)
	if err != nil {
		log.Print(err)
		return "", errors.New("sistemde böyle bir kullanıcı bulunmamaktadır")
	}
	fmt.Println(user)
	user.Type = model.UserTypeAdmin
	err = u.Update(context.Background(), user)
	if err != nil {
		log.Println(err)
		return "", errors.New("bir hata meydana geldi")
	}
	str := fmt.Sprintf("%s kişisi yetkisi admin olarak güncellendi", user.Username)
	return str, nil
}

func (u UserHandler) IsAdmin(id int64) bool { //todo burası tekrar etti. bunu ortak hale getir.
	user, err := u.userService.GetByUserID(context.Background(), id)
	if err != nil {
		log.Println(err)
		return false
	}
	return user.Type == model.UserTypeAdmin
}
