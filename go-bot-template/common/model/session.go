package model

import "github.com/go-bot-template/pkg/model"

type Session struct {
	model.BaseModel
	UserID      int64
	LastCommand string
	State       MessageState
}

type MessageState int

const (
	StateNone            MessageState = 0
	StateWaitingForInput MessageState = 10
)

func (Session) Model() string {
	return ""
}
