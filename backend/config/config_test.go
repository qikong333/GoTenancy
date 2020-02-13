package config

import (
	"regexp"
	"testing"

	"GoTenancy/backend/files"
)

func TestGetAppCreateSysData(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "config",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := GetAppCreateSysData(); got != tt.want {
				t.Errorf("GetAppCreateSysData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAppDriverType(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "config",
			want: "Sqlite",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAppDriverType(); got != tt.want {
				t.Errorf("GetAppDriverType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAppLoggerLevel(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "config",
			want: "debug",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAppLoggerLevel(); got != tt.want {
				t.Errorf("GetAppLoggerLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAppName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "config",
			want: "GoTenancy",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAppName(); got != tt.want {
				t.Errorf("GetAppName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppCookieNameForSessionID(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "config",
			want: "mycookiesessionnameid",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAppCookieNameForSessionID(); got != tt.want {
				t.Errorf("GetAppCookieNameForSessionID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAppURl(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "config",
			want: "irisadminapi.com:80",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAppUrl(); got != tt.want {
				t.Errorf("GetAppURl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMongodbConnect(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "config",
			want: "mongodb://root:123456@127.0.0.1:27017/admin",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMongodbConnect(); got != tt.want {
				t.Errorf("GetMongodbConnect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMysqlConnect(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "config",
			want: "root:([a-z]+)@(127.0.0.1:3306)/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetMysqlConnect()
			matchString, err := regexp.MatchString(tt.want, got)
			if matchString || err != nil {
				t.Errorf("GetMysqlConnect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMysqlName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "config",
			want: "iris",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMysqlName(); got != tt.want {
				t.Errorf("GetMysqlName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMysqlTName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "config",
			want: "tiris",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMysqlTName(); got != tt.want {
				t.Errorf("GetMysqlTName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSqliteConnect(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "config",
			want: files.GetAbsPath("./tmp/gorm.db"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSqliteConnect(); got != tt.want {
				t.Errorf("GetSqliteConnect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSqliteTConnect(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "config",
			want: files.GetAbsPath("./tmp/tgorm.db"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSqliteTConnect(); got != tt.want {
				t.Errorf("GetSqliteTConnect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTestDataName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "config",
			want: "超级管理员",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTestDataName(); got != tt.want {
				t.Errorf("GetTestDataName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTestDataPwd(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "config",
			want: "password",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTestDataPwd(); got != tt.want {
				t.Errorf("GetTestDataPwd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTestDataUserName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "config",
			want: "admin",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTestDataUserName(); got != tt.want {
				t.Errorf("GetTestDataUserName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAdminName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "config",
			want: "超级管理员",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAdminName(); got != tt.want {
				t.Errorf("GetAdminName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAdminPwd(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "config",
			want: "password",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAdminPwd(); got != tt.want {
				t.Errorf("GetAdminPwd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAdminUserName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "config",
			want: "username",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAdminUserName(); got != tt.want {
				t.Errorf("GetAdminUserName() = %v, want %v", got, tt.want)
			}
		})
	}
}
