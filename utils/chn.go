package utils

import (
	"fmt"
	"strings"

	"github.com/spf13/cast"
)

type Number struct{}

// 数字
var chnNumChar = [10]string{"零", "一", "二", "三", "四", "五", "六", "七", "八", "九"}
var chnM = map[string]int{"零": 0, "一": 1, "二": 2, "三": 3, "四": 4, "五": 5, "六": 6, "七": 7, "八": 8, "九": 9}

// 权位
var chnUnitSection = [4]string{"", "万", "亿", "万亿"}

// 数字权位
var chnUnitChar = [4]string{"", "十", "百", "千"}

type chnNameValue struct {
	name    string
	value   int
	secUnit bool
}

// 权位于结点的关系
var chnValuePair = []chnNameValue{
	{"十", 10, false},
	{"百", 100, false},
	{"千", 1000, false},
	{"万", 10000, true},
	{"亿", 100000000, true},
}

// NumberToChinese
// 阿拉伯数字转汉字
func (n Number) NumberToChinese(num int64) (numStr string) {
	var unitPos = 0
	var needZero = false

	for num > 0 { //小于零特殊处理
		section := num % 10000 // 已万为小结处理
		if needZero {
			numStr = chnNumChar[0] + numStr
		}
		strIns := n.sectionToChinese(section)
		if section != 0 {
			strIns += chnUnitSection[unitPos]
		} else {
			strIns += chnUnitSection[0]
		}
		numStr = strIns + numStr
		/*千位是 0 需要在下一个 section 补零*/
		needZero = (section < 1000) && (section > 0)
		num = num / 10000
		unitPos++
	}
	return
}

func (n Number) NumToCn(num interface{}) string {
	tmpNum := cast.ToString(num)
	ss := strings.Split(tmpNum, "")
	tmpCnNum := ""
	for i, v := range ss {
		fmt.Println(i, v)
		tmpCnNum += chnNumChar[cast.ToInt(v)]
	}
	return tmpCnNum
}

func (n Number) sectionToChinese(section int64) (chnStr string) {
	var strIns string
	var unitPos = 0
	var zero = true
	for section > 0 {
		var v = section % 10
		if v == 0 {
			if !zero {
				zero = true /*需要补，zero 的作用是确保对连续的多个，只补一个中文零*/
				chnStr = chnNumChar[v] + chnStr
			}
		} else {
			zero = false                   //至少有一个数字不是
			strIns = chnNumChar[v]         //此位对应的中文数字
			strIns += chnUnitChar[unitPos] //此位对应的中文权位
			chnStr = strIns + chnStr
		}
		unitPos++ //移位
		section = section / 10
	}
	return
}

// ChineseToNumber
// 汉字转阿拉伯数字
func (n Number) ChineseToNumber(chnStr string) (rtnInt int) {
	var section = 0
	var number = 0
	for index, value := range chnStr {

		var num = n.chineseToValue(string(value))
		if num > 0 {
			number = num
			if index == len(chnStr)-3 {
				section += number
				rtnInt += section
				break
			}
		} else {
			unit, secUnit := n.chineseToUnit(string(value))
			if secUnit {
				section = (section + number) * unit
				rtnInt += section
				section = 0

			} else {
				section += number * unit

			}
			number = 0
			if index == len(chnStr)-3 {
				rtnInt += section
				break
			}
		}
	}

	return
}
func (num Number) chineseToUnit(chnStr string) (unit int, secUnit bool) {
	for i := 0; i < len(chnValuePair); i++ {
		if chnValuePair[i].name == chnStr {
			unit = chnValuePair[i].value
			secUnit = chnValuePair[i].secUnit
		}
	}
	return
}

func (num Number) chineseToValue(chnStr string) (rtnNum int) {
	return chnM[chnStr]
}
