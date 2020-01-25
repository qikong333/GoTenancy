package model

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

// NewToken 返回将ID与唯一标识符结合在一起的令牌。
func NewToken(id int64) string {
	return fmt.Sprintf("%d|%s", id, uuid.NewV4().String())
}

// ParseToken 返回给定令牌的id和uuid。
func ParseToken(token string) (int64, string) {
	pairs := strings.Split(token, "|")
	if len(pairs) != 2 {
		return -1, ""
	}

	id, err := strconv.ParseInt(pairs[0], 10, 64)
	if err != nil {
		return -1, ""
	}
	return id, pairs[1]
}

// NewFriendlyID 返回一个唯一的友好ID。
func NewFriendlyID(id int64, key string) string {
	n := time.Now()
	i, _ := strconv.Atoi(
		fmt.Sprintf("%d%d%d%d%d%d%d%d%d",
			id,
			len(key),
			n.Year()-2000,
			int(n.Month()),
			n.Day(),
			n.Hour(),
			n.Minute(),
			n.Second(),
			n.Nanosecond()))
	return fmt.Sprintf("%x", i)
}

func StringToKey(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Printf("error converting %s to int64\n", s)
		return -1
	}
	return i
}
