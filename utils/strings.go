package utils

import (
	"math/rand"
	"strings"
	"time"
)

func FirstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func FirstLower(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// MaheHump 将字符串转换为驼峰命名
func MaheHump(s string) string {
	words := strings.Split(s, "-")

	for i := 1; i < len(words); i++ {
		words[i] = strings.Title(words[i])
	}

	return strings.Join(words, "")
}

// NewTicketId 生成唯一的ticketId
func NewTicketId(prefix string) string {
	fmtTime := time.Now().Format("20060102150405")
	source := "0123456789"
	bytes := []byte(source)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 4; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return prefix + fmtTime + string(result)
}

func RandomString(n int) string {
	charset := "0123456789ABCDEF"
	sb := strings.Builder{}
	sb.Grow(n)
	for i := 0; i < n; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}
