package handlers

import (
	"context"
	userservice "github.com/go-bot-template/app/service"
	"github.com/go-bot-template/pkg/app"
	"github.com/uptrace/bun"
	"log"
)

type ScheduledHandler struct {
	//service     service.IBaseService[model.User]
	userService userservice.IUserService
}

func NewScheduled(db *bun.DB) ScheduledHandler {
	return ScheduledHandler{
		//service:     service.NewBaseService[model.User](db),
		userService: userservice.NewUserService(db),
	}
}

func (u ScheduledHandler) ScheduledDeneme() []app.ScheduledResponse {
	userids, err := u.userService.GetAllUserIDs(context.Background())
	if err != nil {
		//todo bence burada paniklemeli ama bakarız
		log.Println(err)
		return nil
	}

	var response []app.ScheduledResponse

	for _, id := range userids {
		response = append(response, app.ScheduledResponse{
			UserID: id,
			Result: "Saat 01:00'dan önce çetelelerinizi doldurunuz!",
			Error:  nil,
		})
	}
	return response
}
