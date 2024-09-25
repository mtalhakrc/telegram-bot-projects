package model

import "github.com/go-bot-template/pkg/model"

type UserType int

const (
	UserTypeNormal UserType = 0
	UserTypeAdmin  UserType = 10
)

type User struct {
	model.BaseModel
	Name     string
	Username string
	UserID   int64    //telegram user id
	Type     UserType `json:"type"`
}

func (User) Model() string {
	return "user"
}
