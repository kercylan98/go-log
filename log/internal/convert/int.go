package convert

import (
	"github.com/kercylan98/go-log/log/internal/constraints"
	"strconv"
)

// IntToString 将 int 转换为 string 类型
func IntToString(i int) string {
	return strconv.Itoa(i)
}

// Int8ToString 将 int8 转换为 string 类型
func Int8ToString(i int8) string {
	return strconv.Itoa(int(i))
}

// Int16ToString 将 int16 转换为 string 类型
func Int16ToString(i int16) string {
	return strconv.Itoa(int(i))
}

// Int32ToString 将 int32 转换为 string 类型
func Int32ToString(i int32) string {
	return strconv.Itoa(int(i))
}

// Int64ToString 将 int64 转换为 string 类型
func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

// UintToString 将 uint 转换为 string 类型
func UintToString(i uint) string {
	return strconv.FormatUint(uint64(i), 10)
}

// Uint8ToString 将 uint8 转换为 string 类型
func Uint8ToString(i uint8) string {
	return strconv.FormatUint(uint64(i), 10)
}

// Uint16ToString 将 uint16 转换为 string 类型
func Uint16ToString(i uint16) string {
	return strconv.FormatUint(uint64(i), 10)
}

// Uint32ToString 将 uint32 转换为 string 类型
func Uint32ToString(i uint32) string {
	return strconv.FormatUint(uint64(i), 10)
}

// Uint64ToString 将 uint64 转换为 string 类型
func Uint64ToString(i uint64) string {
	return strconv.FormatUint(i, 10)
}

func FastUint64ToString(i uint64) string {
	if i == 0 {
		return "0"
	}

	// 临时缓冲区，uint64 最大长度为 20 个字符
	var buf [20]byte
	pos := len(buf)

	// 使用除法和取模操作将数字转换为字符
	for i >= 10 {
		pos--
		q := i / 10
		buf[pos] = byte('0' + i%10)
		i = q
	}
	pos--
	buf[pos] = byte('0' + i)

	return string(buf[pos:])
}

// IntToBoolean 将 int 转换为 bool 类型
func IntToBoolean[I constraints.Int](i I) bool {
	return i != 0
}
