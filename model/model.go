package model

import (
	"time"
)

// Account 表示一个账号的基本信息
type Account struct {
	ID             int64     `json:"id"`
	Email          string    `json:"email"`
	StripeID       string    `json:"stripeId"`
	SubscriptionID string    `json:"subscriptionId"`
	Plan           string    ` json:"plan"`
	IsYearly       bool      `json:"isYearly"`
	SubscribedOn   time.Time `json:"subscribed"`
	Seats          int       ` json:"seats"`
	TrialInfo      Trial     ` json:"trial"`
	IsActive       bool      ` json:"active"`

	Users []User ` json:"users"`
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
	ID           int64         `json:"id"`
	AccountID    int64         ` json:"accountId"`
	Email        string        `json:"email"`
	Password     string        ` json:"-"`
	Token        string        ` json:"token"`
	Role         Roles         ` json:"role"`
	AccessTokens []AccessToken ` json:"accessTokens"`
}

// AccessToken 表示访问 tokens.
type AccessToken struct {
	ID     int64  ` json:"id"`
	UserID int64  ` json:"userId"`
	Name   string ` json:"name"`
	Token  string ` json:"token"`
}

// APIRequest 表示单个 API 请求。
type APIRequest struct {
	ID         int64     ` json:"id"`
	AccountID  int64     ` json:"accountId"`
	UserID     int64     ` json:"userId"`
	URL        string    `json:"url"`
	Requested  time.Time ` json:"requested"`
	StatusCode int       ` json:"statusCode"`
	RequestID  string    ` json:"reqId"`
}

// Webhook 表示网络订阅。
type Webhook struct {
	ID        int64     `json:"id"`
	AccountID int64     `json:"accountId"`
	EventName string    `json:"event"`
	TargetURL string    `json:"url"`
	IsActive  bool      `json:"active"`
	Created   time.Time `json:"created"`
}
