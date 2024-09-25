package handlers

import (
	"context"
	"errors"
	userservice "github.com/go-bot-template/app/service"
	"github.com/go-bot-template/common/model"
	"github.com/go-bot-template/pkg/app"
	"github.com/go-bot-template/pkg/service"
	"github.com/uptrace/bun"
	"strings"
)

type UserHandler struct {
	service.IBaseService[model.User]
	userService userservice.IUserService
}

func NewUserHandler(db *bun.DB) *UserHandler {
	return &UserHandler{
		IBaseService: service.NewBaseService[model.User](db),
		userService:  userservice.NewUserService(db),
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
		return "", errors.New("isim girmelisin")
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

func (u UserHandler) UpdateName(ctx *app.Ctx, params []string) (string, error) {
	userid := ctx.SentFrom().ID

	//zaten kaydoldu mu kontrol et
	old, err := u.userService.GetByUserID(context.Background(), userid)
	if err != nil {
		err = errors.New("kaydınız bulunamadı")
		return "", err
	}

	newname := strings.Join(params, " ")

	if newname == "" {
		return "", errors.New("isim girmelisin")
	} else if newname == old.Name {
		return "", errors.New("yeni ismin eskisiyle aynı olamaz")
	}

	old.Name = newname

	err = u.Update(context.Background(), old)
	if err != nil {
		return "", err
	}
	return "isminiz başarı ile değiştirildi", nil
}

func (u UserHandler) DeleteUser(ctx *app.Ctx, params []string) (string, error) {
	err := u.DeleteByUserID(context.Background(), ctx.SentFrom().ID)
	return "kaydınız başarı ile silindi", err
}

func (u UserHandler) Deneme(ctx *app.Ctx, params []string) (string, error) {
	return "denem", nil
}
