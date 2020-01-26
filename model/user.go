package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Trial 表示一个账号的使用信息
type Trial struct {
	IsTrial  bool      ` json:"trial"`
	Plan     string    ` json:"plan"`
	Start    time.Time ` json:"start"`
	Extended int       ` json:"extended"`
}

// Roles 和用户访问控制和认证一起使用， 你可以在默认的角色中自定义自己的角色。
type Roles int

const (
	// RolePublic 公共路由访问
	RolePublic Roles = 0
	// RoleFree 免费用户
	RoleFree = 10
	// RoleUser 标准用户
	RoleUser = 20
	// RoleAdmin 管理员
	RoleAdmin = 99
)

// User 表示一个用户
type User struct {
	gorm.Model

	AccountID    uint
	Email        string `gorm:"type:varchar(255);column:email"`
	Password     string `gorm:"type:varchar(255);column:-"`
	Token        string `gorm:"type:varchar(255);column:token"`
	Role         Roles
	AccessTokens []AccessToken
}
