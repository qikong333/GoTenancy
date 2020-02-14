package session

import (
	"testing"
)

func TestSingle(t *testing.T) {

	if got := Singleton(); got == nil {
		t.Errorf("Single() = %v, is %v", got, nil)
	}

}
