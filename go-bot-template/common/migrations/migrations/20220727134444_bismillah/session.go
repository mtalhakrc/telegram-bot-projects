package migration

import "github.com/go-bot-template/pkg/model"

type Session struct {
	model.BaseModel
	UserID      int64
	LastCommand string
	Status      MessageStatus
}

type MessageStatus int

const (
	New             MessageStatus = 0
	WaitingForInput MessageStatus = 10
)

func (Session) Model() string {
	return ""
}
