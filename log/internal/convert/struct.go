package convert

import (
	"encoding/json"
	jsonIter "github.com/json-iterator/go"
	"unsafe"
)

// PointerTo 将一个结构体指针转换为另一个结构体指针
//   - 两个结构体字段必须完全一致
//   - 该函数可以绕过私有字段的访问限制
func PointerTo[S, D any](src *S) (dst *D) {
	return (*D)(unsafe.Pointer(src))
}

// StructToWithJSON 将一个结构体转换为另一个结构体
//   - 仅支持满足 JSON 序列化的结构体字段
func StructToWithJSON[S, D any](src *S) (dst *D) {
	dst = new(D)
	var jsonData, _ = jsonIter.ConfigCompatibleWithStandardLibrary.Marshal(src)
	_ = json.Unmarshal(jsonData, dst)
	return
}
