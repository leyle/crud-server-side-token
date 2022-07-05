package sst

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func encodeUserId(userId string) string {
	t := time.Now().Unix()
	return fmt.Sprintf("%s|%d", userId, t)
}

func decodeUserId(val string) (string, int64, error) {
	if !strings.Contains(val, "|") {
		return "", 0, fmt.Errorf("invalid encodedUserId[%s] format, no `|` in text", val)
	}

	result := strings.Split(val, "|")
	userId := result[0]
	strT := result[1]
	t, err := strconv.ParseInt(strT, 10, 64)
	if err != nil {
		return "", 0, err
	}

	return userId, t, nil
}
