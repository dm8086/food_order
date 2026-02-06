package utils

import (
	"github.com/mozillazg/go-pinyin"
)

// GetFirstLetter 获取首字母
func GetFirstLetter(str string) string {
	a := pinyin.NewArgs()
	a.Style = pinyin.FirstLetter
	letterArr := pinyin.Pinyin(str, a)
	var letter string
	for _, v := range letterArr {
		letter += v[0]
	}
	return letter
}

// IsLetter 判断是否是字母
func IsLetter(str string) bool {
	for _, v := range str {
		if (v < 'a' || v > 'z') && (v < 'A' || v > 'Z') {
			return false
		}
	}
	return true
}
