package utils

import (
	"bytes"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
	"unsafe"

	"github.com/mozillazg/go-pinyin"
	"github.com/spf13/cast"
)

var (
	camelRe = regexp.MustCompile("(_)([a-zA-Z]+)")
	snakeRe = regexp.MustCompile("([a-z0-9])([A-Z])")
)

type Str struct{}

func (s Str) StringToByteSlice(str string) []byte {
	tmp1 := (*[2]uintptr)(unsafe.Pointer(&str))
	tmp2 := [3]uintptr{tmp1[0], tmp1[1], tmp1[1]}
	return *(*[]byte)(unsafe.Pointer(&tmp2))

}

func (s Str) ByteSliceToString(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}

// StrPad
// input string 原字符串
// padLength int 规定补齐后的字符串位数
// padString string 自定义填充字符串
// padType string 填充类型:LEFT(向左填充,自动补齐位数), 默认右侧
func (s Str) StrPad(input string, padString string, padLength int, padType string) string {
	output := ""
	inputLen := len(input)

	if inputLen >= padLength {
		return input
	}
	padStringLen := len(padString)
	needFillLen := padLength - inputLen

	if diffLen := padStringLen - needFillLen; diffLen > 0 {
		padString = padString[diffLen:]
	}
	for i := 1; i <= needFillLen; i += padStringLen {
		output += padString
	}
	switch padType {
	case "LEFT":
		return output + input
	default:
		return input + output
	}
}

// IsEmptyPtr
// 判断字符串指针是否为空
func (s Str) IsEmptyPtr(inStr *string) bool {
	return inStr == nil || s.IsNullStr(*inStr)
}

// IsNullStr
// 判断字符串是否为空 "" / "null" / "nil" / "undefined"
func (s Str) IsNullStr(inStr string) bool {
	inStr = strings.TrimSpace(inStr)
	return inStr == "" || inStr == "null" || inStr == "nil" || inStr == "undefined"
}

// StringUnquote 去除字符串两边的双引号
func (s Str) StringUnquote(value string) string {
	if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) {
		value, _ = strconv.Unquote(value)
	}
	return value
}

// AddStr 组装字符串
func (s Str) AddStr(args ...interface{}) string {
	var buffer bytes.Buffer //Buffer是一个实现了读写方法的可变大小的字节缓冲
	for _, arg := range args {
		buffer.WriteString(cast.ToString(arg))
	}
	return buffer.String()
}

// AddStrEx
// 组装字符串, 添加分割字符串
func (s Str) AddStrEx(split string, args ...interface{}) string {
	var buffer bytes.Buffer //Buffer是一个实现了读写方法的可变大小的字节缓冲
	for _, arg := range args {
		buffer.WriteString(cast.ToString(arg))
		buffer.WriteString(split)
	}
	return buffer.String()
}

// GetPYM 获取拼音码
func (s Str) GetPYM(str string) string {
	tmpStr := ""
	data := pinyin.LazyConvert(str, nil)
	for _, v := range data {
		tmpStr += v[0:1]
	}
	return strings.ToUpper(tmpStr)
}

// GetRandomString
// 随机生成一个指定位数的字符串
func (s Str) GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// StringBuilder 高性能构建字符串工具函数
func (s Str) StringBuilder(values ...interface{}) string {
	builder := strings.Builder{}
	for _, value := range values {
		_, _ = fmt.Fprintf(&builder, "%v", value)
	}
	return builder.String()
}

// FillZero
// 补 0 返回字符串
//
//	value：传入值；length：总位数
func (s Str) FillZero(value, length int) string {
	if length <= 0 {
		length = len(strconv.Itoa(value))
	}
	strVal := strconv.Itoa(value)
	count := length - len(strVal)
	if count > 0 {
		strVal = strings.Repeat("0", count) + strVal
	}
	return strVal
}

// IsStrHasAnyPrefix 判断字符串 s 是否以 prefixes 中的任意一个前缀开头
func (s Str) IsStrHasAnyPrefix(str string, prefixes []string) (prefixIndex int, has bool) {
	for i, prefix := range prefixes {
		if has = strings.HasPrefix(str, prefix); has {
			prefixIndex = i
			break
		}
	}
	return
}

// StrAllLetter
// 判断字符串是否全部是字母
func (s Str) StrAllLetter(str string) bool {
	match, _ := regexp.MatchString(`^[A-Za-z]+$`, str)
	return match
}

// SubStr
// 截取指定长度的字符串
func (s Str) SubStr(str string, length int) (ret string) {
	var count int
	for i := 0; i < len(str); {
		r, size := utf8.DecodeRuneInString(str[i:])
		i += size
		ret += string(r)
		count++
		if length <= count {
			break
		}
	}
	return
}
