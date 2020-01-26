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
func NewToken(id uint) string {
	return fmt.Sprintf("%d|%s", id, uuid.NewV4().String())
}

// ParseToken 返回给定令牌的id和uuid。
func ParseToken(token string) (uint, string) {
	pairs := strings.Split(token, "|")
	if len(pairs) != 2 {
		return 0, ""
	}

	id, err := strconv.ParseUint(pairs[0], 10, 64)
	if err != nil {
		return 0, ""
	}
	return uint(id), pairs[1]
}

// NewFriendlyID 返回一个唯一的友好ID。
func NewFriendlyID(id uint, key string) string {
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

func StringToKey(s string) uint {
	id, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		log.Printf("error converting %s to uint\n", s)
		return 0
	}
	return uint(id)
}
