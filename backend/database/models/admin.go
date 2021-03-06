package models

import (
	"fmt"
	"time"

	"GoTenancy/backend/database"
	"GoTenancy/backend/libs"
	"GoTenancy/backend/validates"
	"github.com/fatih/color"
	"github.com/jinzhu/gorm"
)

type Admin struct {
	gorm.Model

	Name     string `gorm:"not null VARCHAR(191)"`
	Username string `gorm:"unique;VARCHAR(191)"`
	Password string `gorm:"not null VARCHAR(191)"`
}

func NewAdmin(id uint, username string) *Admin {
	return &Admin{
		Model: gorm.Model{
			ID:        id,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Username: username,
	}
}

func NewAdminByStruct(ru *validates.CreateUpdateAdminRequest) *Admin {
	password, err := libs.GeneratePassword(ru.Password)
	if err != nil {
		color.Red(fmt.Sprintf("NewAdminByStruct:%s \n ", err))
		return nil
	}

	return &Admin{
		Model: gorm.Model{
			ID:        0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Username: ru.Username,
		Name:     ru.Name,
		Password: password,
	}
}

func (u *Admin) GetAdminByUserName() {
	IsNotFound(database.GetGdb().Where("username = ?", u.Username).First(u).Error)
}

func (u *Admin) GetAdminById() {
	IsNotFound(database.GetGdb().Where("id = ?", u.ID).First(u).Error)
}

/**
 * 通过 id 删除用户
 * @method DeleteAdminById
 */
func (u *Admin) DeleteAdmin() {
	if err := database.GetGdb().Delete(u).Error; err != nil {
		color.Red(fmt.Sprintf("DeleteAdminByIdErr:%s \n ", err))
	}
}

/**
 * 获取所有的账号
 * @method GetAllAdmin
 * @param  {[type]} name string [description]
 * @param  {[type]} username string [description]
 * @param  {[type]} orderBy string [description]
 * @param  {[type]} offset int    [description]
 * @param  {[type]} limit int    [description]
 */
func GetAllAdmins(name, orderBy string, offset, limit int) []*Admin {
	var users []*Admin
	q := GetAll(name, orderBy, offset, limit)
	if err := q.Find(&users).Error; err != nil {
		color.Red(fmt.Sprintf("GetAllAdminErr:%s \n ", err))
		return nil
	}
	return users
}

/**
 * 创建
 * @method CreateAdmin
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func (u *Admin) CreateAdmin(aul *validates.CreateUpdateAdminRequest) {
	password, err := libs.GeneratePassword(aul.Password)
	if err != nil {
		color.Red(fmt.Sprintf("CreateAdminErr:%s \n ", err))
	}
	u.Password = password
	if err := database.GetGdb().Create(u).Error; err != nil {
		color.Red(fmt.Sprintf("CreateAdminErr:%s \n ", err))
	}

	return
}

/**
 * 更新
 * @method UpdateAdmin
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func (u *Admin) UpdateAdmin(uj *validates.CreateUpdateAdminRequest) {
	password, err := libs.GeneratePassword(uj.Password)
	if err != nil {
		color.Red(fmt.Sprintf("UpdateAdminErr:%s \n ", err))
	}
	uj.Password = password
	if err := Update(u, uj); err != nil {
		color.Red(fmt.Sprintf("UpdateAdminErr:%s \n ", err))
	}
}

/**
 * 判断用户是否登录
 * @method CheckLogin
 * @param  {[type]}  id       int    [description]
 * @param  {[type]}  password string [description]
 */
func (u *Admin) CheckLogin(password string) (bool, string) {
	if u.ID == 0 {
		return false, "用户不存在"
	} else {
		if ok, _ := libs.ValidatePassword(password, u.Password); ok {
			return true, "登陆成功"
		} else {
			return false, "用户名或密码错误"
		}
	}
}
