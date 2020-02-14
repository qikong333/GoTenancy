package redis

import (
	"testing"
)

func TestSingleton(t *testing.T) {

	if got := Singleton(); got == nil {
		t.Errorf("Singleton() = %v, is  %v", got, nil)
	}

}
