package model

import "time"

type BaseModel struct {
	ID        int64     `json:"id" bun:"id,pk,autoincrement"`
	CreatedAt time.Time `json:"created_at" bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at" bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt time.Time `json:"deleted_at" bun:",soft_delete,nullzero"`
}

type IModel interface {
	Model() string
}
