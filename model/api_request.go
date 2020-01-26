package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// APIRequest 表示单个 API 请求。
type APIRequest struct {
	gorm.Model

	AccountID  uint
	UserID     uint
	URL        string `gorm:"type:varchar(255);column:url"`
	Requested  time.Time
	StatusCode int
	RequestID  string `gorm:"type:varchar(255);column:reqId"`
}
