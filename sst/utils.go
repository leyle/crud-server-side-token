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

	splitNum := 2
	result := strings.Split(val, "|")
	if len(result) != splitNum {
		return "", 0, fmt.Errorf("invalid encodedUserId[%s] format, more than one `|` in text", val)
	}
	userId := result[0]
	strT := result[1]
	t, err := strconv.ParseInt(strT, 10, 64)
	if err != nil {
		return "", 0, err
	}

	return userId, t, nil
}
