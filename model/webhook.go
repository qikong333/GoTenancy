package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Webhook 表示网络订阅。
type Webhook struct {
	gorm.Model

	AccountID uint
	EventName string `gorm:"type:varchar(255);column:event"`
	TargetURL string `gorm:"type:varchar(255);column:url"`
	IsActive  bool
	Created   time.Time
}
