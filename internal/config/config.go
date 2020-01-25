package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/snowlyg/GoTenancy/data"
)

// EmailProvider 定义电子邮件提供商的可能实现。
type EmailProvider string

const (
	// EmailProviderSES 使用 Amazon SES 邮件服务。
	EmailProviderSES EmailProvider = "amazonses"
)

// Configuration 定义整个库使用的配置
type Configuration struct {
	EmailFrom     string        `json:"emailFrom"`
	EmailFromName string        `json:"emailFromName"`
	EmailProvider EmailProvider `json:"emailProvider"`

	StripeKey string             `json:"stripeKey"`
	Plans     []data.BillingPlan `json:"plans"`

	SignUpTemplate            string `json:"signupTemplate"`
	SignUpSendEmailValidation bool   `json:"sendEmailValidation"`
	SignUpSuccessRedirect     string `json:"signupSuccessRedirect"`
	SignUpErrorRedirect       string `json:"signupErrorRedirect"`
	SignInTemplate            string `json:"signinTemplate"`
	SignInSuccessRedirect     string `json:"signinSuccessRedirect"`
	SignInErrorRedirect       string `json:"signinErrorRedirect"`
}

// Current 保留当前配置
var Current Configuration

// LoadFromFile 加载 ./GoTenancy.json 文件作为默认的库配置
func LoadFromFile() error {
	b, err := ioutil.ReadFile("./GoTenancy.json")
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, &Current); err != nil {
		return fmt.Errorf("error parsing your GoTenancy.json config file: %v", err)
	}
	return nil
}

// Configure 为库的各个重要方面设置适当的值来
// 控制重要程序的行为，比如邮件和注册。
//
// 创建一个 "GoTenancy.json" 配置文件，它将会在启动时自动加载。
func Configure(conf Configuration) {
	Current = conf
}
