package cache

import (
	"testing"
)

func Test_Request_Log(t *testing.T) {
	t.Parallel()

	fakeReq := []byte("ok this is\nthe\nrequest\tlogged")
	if err := LogWebRequest("testing", fakeReq); err != nil {
		t.Fatal(err)
	}

	// 我们希望最后一个插入，所以我们可以测试
	reqID, b, err := GetWebRequest(false)
	if err != nil {
		t.Error(err)
	} else if reqID != "testing" {
		t.Errorf("reqID was %s and we were expecting testing", reqID)
	} else if string(b) != string(fakeReq) {
		t.Errorf("the returned request was %s and we were looking for %s", string(b), string(fakeReq))
	}
}
