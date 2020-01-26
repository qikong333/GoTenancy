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
