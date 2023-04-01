package utils

import (
	"math"
	"reflect"
	"strings"
)

type Slice struct{}

// DeleteElement
// 字符串切片删除元素
func (s Slice) DeleteElement(slice []string, element string) []string {
	for i, v := range slice {
		if v == element {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

// InArray
// 判断元素是否在数组中
func (s Slice) InArray(arr []string, element string) bool {
	for _, v := range arr {
		if v == element {
			return true
		}
	}
	return false
}

func (s Slice) ListToMap(list interface{}, key string) map[string]interface{} {
	res := make(map[string]interface{})
	arr := s.ToSlice(list)
	for _, row := range arr {
		immutable := reflect.ValueOf(row)
		var val string
		if immutable.Kind() == reflect.Ptr {
			val = immutable.Elem().FieldByName(key).String()
		} else {
			val = immutable.FieldByName(key).String()
		}
		res[val] = row
	}
	return res
}

func (s Slice) ToSlice(arr interface{}) []interface{} {
	ret := make([]interface{}, 0)
	v := reflect.ValueOf(arr)
	if v.Kind() != reflect.Slice {
		ret = append(ret, arr)
		return ret
	}
	l := v.Len()
	for i := 0; i < l; i++ {
		ret = append(ret, v.Index(i).Interface())
	}
	return ret
}

// FilterNoEmptyRepeatValues
// 过滤切片数组中非空且不重复的值
func (s Slice) FilterNoEmptyRepeatValues(values []string) (filterValues []string) {
	for _, value := range values {
		value = strings.TrimSpace(value)
		_, existed := s.SliceContainsStr(filterValues, value)
		if existed || value == "" {
			continue
		}
		filterValues = append(filterValues, value)
	}
	return
}

// SliceContainsStr 判断 string 切片中是否包含某个值
func (s Slice) SliceContainsStr(slice []string, value string) (index int, has bool) {
	for i, s := range slice {
		if s == value {
			index = i
			has = true
			break
		}
	}
	return
}

// SliceContainsInt 判断 int 切片中是否包含某个值
func (s Slice) SliceContainsInt(slice []int, value int) (index int, has bool) {
	for i, s := range slice {
		if s == value {
			index = i
			has = true
			break
		}
	}
	return
}

// SliceRemoveRepeatedStr 字符串型切片去重
func (s Slice) SliceRemoveRepeatedStr(strings []string) []string {
	result := make([]string, 0)
	// map 用于保存已存在的元素
	m := make(map[string]bool)
	for _, v := range strings {
		if _, ok := m[v]; !ok {
			result = append(result, v)
			m[v] = true
		}
	}
	return result
}

// SliceRemoveRepeatedInt 整型切片去重
func (s Slice) SliceRemoveRepeatedInt(integers []int) []int {
	result := make([]int, 0)
	// map 用于保存已存在的元素
	m := make(map[int]bool)
	for _, v := range integers {
		if _, ok := m[v]; !ok {
			result = append(result, v)
			m[v] = true
		}
	}
	return result
}

// SliceRemoveStr 移除字符串型切片中的指定元素
func (s Slice) SliceRemoveStr(strings []string, str string) []string {
	for i := 0; i < len(strings); i++ {
		if str == strings[i] {
			strings = append(strings[:i], strings[i+1:]...)
			i = i - 1
		}
	}
	return strings
}

// SliceRemoveInt 移除字符串型切片中的指定元素
func (s Slice) SliceRemoveInt(integers []int, integer int) []int {
	for i := 0; i < len(integers); i++ {
		if integer == integers[i] {
			integers = append(integers[:i], integers[i+1:]...)
			i = i - 1
		}
	}
	return integers
}

// SlicePage 切片分页
func (s Slice) SlicePage(page, pageSize, nums int) (sliceStart, sliceEnd int) {
	if page < 0 {
		page = 1
	}

	if pageSize < 0 {
		pageSize = 10
	}

	if pageSize > nums {
		return 0, nums
	}

	// 总页数
	pageCount := int(math.Ceil(float64(nums) / float64(pageSize)))
	if page > pageCount {
		return 0, 0
	}
	sliceStart = (page - 1) * pageSize
	sliceEnd = sliceStart + pageSize

	if sliceEnd > nums {
		sliceEnd = nums
	}
	return sliceStart, sliceEnd
}
