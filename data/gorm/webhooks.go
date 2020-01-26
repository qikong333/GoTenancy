package gorm

import (
	"errors"

	"github.com/snowlyg/GoTenancy/data"
	"github.com/snowlyg/GoTenancy/model"
)

type Webhooks struct {
	DB *data.DB
}

func (wh *Webhooks) Add(accountID uint, events, url string) error {
	webhook := model.Webhook{
		AccountID: accountID,
		EventName: events,
		TargetURL: url,
	}
	if wh.DB.Connection.NewRecord(&webhook) {
		return errors.New("model.Account 创建失败")
	}

	return nil

}

func (wh *Webhooks) List(accountID uint) ([]model.Webhook, error) {
	var hooks []model.Webhook
	if err := wh.DB.Connection.Where("account_id", accountID).Find(&hooks).Error; err != nil {
		return nil, err
	}

	return hooks, nil
}

func (wh *Webhooks) Delete(accountID uint, event, url string) error {

	if err := wh.DB.Connection.Where("account_id = ?", accountID).
		Where("events = ?", event).
		Where("url = ?", url).Delete(wh).Error; err != nil {
		return err
	}
	return nil
}

func (wh *Webhooks) AllSubscriptions(event string) ([]model.Webhook, error) {
	var hooks []model.Webhook
	if err := wh.DB.Connection.Where("events", event).Find(&hooks).Error; err != nil {
		return nil, err
	}

	return hooks, nil
}
