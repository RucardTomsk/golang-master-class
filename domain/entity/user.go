package entity

import (
	"UGUAPI/domain/base"
	"UGUAPI/domain/enum"
	"github.com/google/uuid"
)

type User struct {
	base.EntityWithIdKey
	Name     string
	Email    string `gorm:"uniqueIndex"`
	Password string
	Role     enum.TypeRole
	Sessions []Session
}

type Session struct {
	base.EntityWithIdKey
	User   *User
	UserID uuid.UUID
}
