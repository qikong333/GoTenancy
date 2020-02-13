package session

import (
	"reflect"
	"testing"

	"GoTenancy/backend/config"
	"github.com/kataras/iris/v12/sessions"
)

func TestSingle(t *testing.T) {
	tests := []struct {
		name string
		want *sessions.Sessions
	}{
		{
			name: "test-session",
			want: sessions.New(sessions.Config{Cookie: config.GetAppCookieNameForSessionID(), AllowReclaim: true}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Single(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Single() = %v, want %v", got, tt.want)
			}
		})
	}
}
