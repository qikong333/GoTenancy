package data

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/snowlyg/GoTenancy/model"
)

// DB 是一个包含数据库连接引用的数据库抽象接口
// 目前，支持Postgres和内存数据提供程序。
type DB struct {
	// DatabaseName 数据库名称
	DatabaseName string
	// Connection据库连接的引用。
	Connection *gorm.DB

	// Users 包含账号, 用户和开票相关访问方法的数据
	Users UserServices
	// Webhooks 包含管理 Webhooks 相关访问方法的数据。
	Webhooks WebhookServices
}

// UserServices 是一个包含账号, 用户和开票相关所有方法的接口
type UserServices interface {
	SignUp(email, password string) (*model.Account, error)
	ChangePassword(id, accountID int64, passwd string) error
	AddToken(accountID, userID int64, name string) (*model.AccessToken, error)
	RemoveToken(accountID, userID, tokenID int64) error
	Auth(accountID int64, token string, pat bool) (*model.Account, *model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	GetDetail(id int64) (*model.Account, error)
	GetByStripe(stripeID string) (*model.Account, error)
	SetSeats(id int64, seats int) error
	ConvertToPaid(id int64, stripeID, subID, plan string, yearly bool, seats int) error
	ChangePlan(id int64, plan string, yearly bool) error
	Cancel(id int64) error
}

// AdminServices TODO: investigate this...
type AdminServices interface {
	LogRequests(reqs []model.APIRequest) error
}

// WebhookServices 一个包含管理 webhook 所有方法的接口
type WebhookServices interface {
	Add(accountID int64, events, url string) error
	List(accountID int64) ([]model.Webhook, error)
	Delete(accountID int64, event, url string) error
	AllSubscriptions(event string) ([]model.Webhook, error)
}

// NewID 根据帐户和用户ID返回每秒的唯一字符串。
func NewID(accountID, userID int64) string {
	n := time.Now()
	i, _ := strconv.Atoi(
		fmt.Sprintf("%d%d%d%d%d%d%d%d",
			accountID,
			userID,
			n.Year()-2000,
			int(n.Month()),
			n.Day(),
			n.Hour(),
			n.Minute(),
			n.Second()))
	return fmt.Sprintf("%x", i)
}
