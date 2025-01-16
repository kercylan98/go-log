package charproc

import (
	"strings"
	"unicode"
)

// FirstUpper 返回字符串的首字母大写形式
func FirstUpper(s string) string {
	if len(s) == 0 {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

// BigCamel 大驼峰字符串
func BigCamel(str string) string {
	return FirstUpper(Camel(str))
}

// Camel 驼峰字符串
func Camel(str string) string {
	var builder strings.Builder
	vv := []rune(str)
	for i := 0; i < len(vv); i++ {
		if vv[i] == '_' {
			i++
			if vv[i] >= 97 && vv[i] <= 122 {
				vv[i] -= 32
				builder.WriteRune(vv[i])
			} else {
				return str
			}
		} else {
			builder.WriteRune(vv[i])
		}
	}
	return builder.String()
}
