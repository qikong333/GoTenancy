package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Account 表示一个账号的基本信息
type Account struct {
	gorm.Model

	Email          string `gorm:"type:varchar(255);column:email"`
	StripeID       string `gorm:"type:varchar(255);column:stripeId"`
	SubscriptionID string `gorm:"type:varchar(255);column:subscriptionId"`
	Plan           string `gorm:"type:varchar(255);column:plan"`
	IsYearly       bool
	SubscribedOn   time.Time `gorm:"column:subscribed"`
	Seats          int
	TrialInfo      Trial `gorm:"type:longText;column:trial"`
	IsActive       bool

	Users []User
}

// IsPaid 账号是否是付费用户
func (a *Account) IsPaid() bool {
	return len(a.StripeID) > 0 && len(a.SubscriptionID) > 0
}

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

	AccountID    int64
	Email        string `gorm:"type:varchar(255);column:email"`
	Password     string `gorm:"type:varchar(255);column:-"`
	Token        string `gorm:"type:varchar(255);column:token"`
	Role         Roles
	AccessTokens []AccessToken
}

// AccessToken 表示访问 tokens.
type AccessToken struct {
	gorm.Model

	UserID int64
	Name   string `gorm:"type:varchar(255);column:name"`
	Token  string `gorm:"type:varchar(255);column:token"`
}

// APIRequest 表示单个 API 请求。
type APIRequest struct {
	gorm.Model

	AccountID  int64
	UserID     int64
	URL        string `gorm:"type:varchar(255);column:url"`
	Requested  time.Time
	StatusCode int
	RequestID  string `gorm:"type:varchar(255);column:reqId"`
}

// Webhook 表示网络订阅。
type Webhook struct {
	gorm.Model

	AccountID int64
	EventName string `gorm:"type:varchar(255);column:event"`
	TargetURL string `gorm:"type:varchar(255);column:url"`
	IsActive  bool
	Created   time.Time
}
