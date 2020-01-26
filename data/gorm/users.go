package gorm

import (
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/snowlyg/GoTenancy/data"
	"github.com/snowlyg/GoTenancy/model"
)

type Users struct {
	DB *data.DB
}

func (u *Users) SignUp(email, password string) (*model.Account, error) {
	account := model.Account{
		Email:          email,
		StripeID:       "",
		SubscriptionID: "",
		Plan:           "",
		IsYearly:       false,
		SubscribedOn:   time.Now(),
		Seats:          0,
		IsActive:       true,
	}

	if err := u.DB.Connection.Create(&account).Error; err != nil {
		return nil, errors.New(fmt.Sprintf("create model.Account error:%v", err))
	}

	token := model.NewToken(account.ID)
	user := model.User{
		AccountID: account.ID,
		Email:     email,
		Password:  password,
		Token:     token,
		Role:      model.RoleAdmin,
	}

	if err := u.DB.Connection.Create(&user).Error; err != nil {
		return nil, errors.New("create model.User error")
	}

	return u.GetDetail(account.ID)
}

func (u *Users) Auth(accountID uint, token string, pat bool) (*model.Account, *model.User, error) {
	token = fmt.Sprintf("%d|%s", accountID, token)

	user := &model.User{AccountID: accountID, Token: token}

	if err := u.DB.Connection.First(user).Error; err != nil {
		return nil, nil, err
	}

	account, err := u.GetDetail(user.AccountID)
	if err != nil {
		return nil, nil, errors.New("user Auth error")
	}

	return account, user, nil
}

func (u *Users) GetDetail(id uint) (*model.Account, error) {
	account := &model.Account{Model: gorm.Model{ID: id}}
	if err := u.DB.Connection.Preload("Users").First(account).Error; err != nil {
		return nil, errors.New("GetDetail error")
	}

	return account, nil
}

func (u *Users) GetUserByEmail(email string) (*model.User, error) {
	user := &model.User{Email: email}
	if err := u.DB.Connection.First(user).Error; err != nil {
		return nil, errors.New("GetUserByEmail error")
	}

	return user, nil
}

func (u *Users) GetByStripe(stripeID string) (*model.Account, error) {
	account := model.Account{StripeID: stripeID}
	if err := u.DB.Connection.First(account).Error; err != nil {
		return nil, errors.New("GetByStripe error")
	}

	return u.GetDetail(account.ID)
}

func (u *Users) ChangePassword(id, accountID uint, passwd string) error {
	if err := u.DB.Connection.Model(&u).
		Where("id = ?", id).
		Where("account_id = ?", accountID).
		Update("password", passwd).Error; err != nil {
		return errors.New("ChangePassword error")
	}

	return nil
}

func (u *Users) SetSeats(id uint, seats int) error {
	if err := u.DB.Connection.Model(&u).
		Where("id = ?", id).
		Update("seats", seats).Error; err != nil {
		return errors.New("SetSeats error")
	}

	return nil
}

func (u *Users) ConvertToPaid(id uint, stripeID, subID, plan string, yearly bool, seats int) error {
	d := map[string]interface{}{
		"stripe_id":       stripeID,
		"subscription_id": subID,
		"subscribed_on":   time.Now(),
		"plan":            plan,
		"seats":           seats,
		"is_yearly":       yearly,
	}

	if err := u.DB.Connection.Model(&u).
		Where("id = ?", id).
		Updates(d).Error; err != nil {
		return errors.New("ConvertToPaid error")
	}

	return nil
}

func (u *Users) ChangePlan(id uint, plan string, yearly bool) error {
	d := map[string]interface{}{
		"plan":      plan,
		"is_yearly": yearly,
	}

	if err := u.DB.Connection.Model(&u).
		Where("id = ?", id).
		Updates(d).Error; err != nil {
		return errors.New("ChangePlan error")
	}

	return nil
}

func (u *Users) Cancel(id uint) error {
	d := map[string]interface{}{
		"subscription_id": "",
		"plan":            "",
		"is_yearly":       false,
	}

	if err := u.DB.Connection.Model(&u).
		Where("id = ?", id).
		Updates(d).Error; err != nil {
		return errors.New("user Cancel error")
	}

	return nil
}

func (u *Users) AddToken(accountID, userID uint, name string) (*model.AccessToken, error) {
	return nil, fmt.Errorf("not implemented")
}

func (u *Users) RemoveToken(accountID, userID, tokenID uint) error {
	return fmt.Errorf("not implemented")
}
