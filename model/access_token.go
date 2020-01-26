package model

import (
	"github.com/jinzhu/gorm"
)

// AccessToken 表示访问 tokens.
type AccessToken struct {
	gorm.Model

	UserID uint
	Name   string `gorm:"type:varchar(255);column:name"`
	Token  string `gorm:"type:varchar(255);column:token"`
}
