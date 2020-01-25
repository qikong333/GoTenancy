package cache

import (
	"bytes"
	"encoding/gob"
)

// LogRequest 将新项目添加到要记录的挂起请求列表中。
func LogRequest(v interface{}) error {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(v); err != nil {
		return err
	}

	if _, err := rc.RPush("reqlog", buf.String()).Result(); err != nil {
		return err
	}
	return nil
}

// DequeueRequests 返回所有准备插入到数据库中的挂起请求。
func DequeueRequests() ([]string, error) {
	return rc.LRange("reqlog", 0, -1).Result()
}
