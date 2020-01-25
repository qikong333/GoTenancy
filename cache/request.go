package cache

import (
	"strings"

	"fmt"

	"github.com/go-redis/redis"
)

// CountWebRequest 返回待分析的失败 Web 请求数。
func CountWebRequest() (int64, error) {
	return rc.LLen("reqs").Result()
}

// GetWebRequest 返回从列表中记录的下一个 Web 请求。
func GetWebRequest(first bool) (reqID string, b []byte, err error) {
	var s string
	if first {
		s, err = rc.LPop("reqs").Result()
	} else {
		s, err = rc.RPop("reqs").Result()
	}

	if err != nil {
		if err == redis.Nil {
			err = nil
		}
		return
	}

	buf := strings.Split(s, "\n|\n")
	if len(buf) != 2 {
		return reqID, b, fmt.Errorf("unable to split request result")
	}

	b = []byte(buf[0])
	reqID = buf[1]

	return
}

// LogWebRequest 保存 Web 请求以进行进一步分析。
func LogWebRequest(reqID string, b []byte) error {
	r := []byte(fmt.Sprintf("\n|\n%s", reqID))
	b = append(b, r...)
	return rc.RPush("reqs", b).Err()
}
