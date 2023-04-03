package utils

import (
	"bytes"
	"database/sql/driver"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-module/carbon/v2"
)

const (
	/* 各个格式化时间模版 */

	FormatDateTimeNum = "20060102150405"          // 日期时间格式数字串精确到秒
	FormatDateTimeMs  = "2006-01-02 15:04:05.999" // 日期时间格式精确到毫秒
	FormatDateTimeSec = "2006-01-02 15:04:05"     // 日期时间格式精确到秒
	FormatDateTime    = "2006-01-02 15:04"        // 日期时间格式精确到分
	FormatDate        = "2006-01-02"              // 日期格式：年-月-日，月日补 0
	FormatDateShort   = "2006-1-2"                // 日期格式：年-月-日
	FormatDateNum     = "20060102"                // 日期格式数字串：年月日
	FormatTimeMs      = "15:04:05.999"            // 时间格式精确到毫秒
	FormatTimeSec     = "15:04:05"                // 时间格式精确到秒
	FormatTime        = "15:04"                   // 时间格式精确到分
	FormatYear        = "2006"                    // 日期年份
	FormatMonth       = "01"
	FormatDay         = "02"
	FormatHour        = "15"
)

type Ltime struct{}

// ParseHSM
// 解析时分秒字符串
func (l Ltime) ParseHSM(timeStr string) (time.Time, error) {
	now := time.Now()
	tmpTimeStr := ""
	if len(timeStr) <= 5 {
		tmpTimeStr = now.Format(FormatDate) + " " + timeStr + ":00"
	}
	if len(timeStr) > 5 && len(timeStr) <= 8 {
		tmpTimeStr = now.Format(FormatDate) + " " + timeStr
	}
	return time.ParseInLocation(FormatDateTimeSec, tmpTimeStr, time.Local)
}

func (l Ltime) ParseDate(timeStr string) time.Time {
	t, err := time.Parse(FormatDate, timeStr)
	if err != nil {
		panic(err)
	}
	return t
}

// InitTime
// 初始化时间
func (l Ltime) InitTime() time.Time {
	t, _ := time.Parse(FormatDateTimeMs, "1899-12-30 23:59:59.999")
	return t
}

func (l Ltime) GetFormatDatetime(t time.Time) string {
	return t.Format(FormatDateTimeMs)
}

// DayDiff
// 日期差
func (l Ltime) DayDiff(t1, t2 time.Time) int64 {
	return carbon.Parse(t1.Format(FormatDateTime)).DiffAbsInDays(carbon.Parse(t2.Format(FormatDateTime)))
}

func (l Ltime) GetNowStr() string {
	return time.Now().Format(FormatDateTimeMs)
}

// NowTimeStr
// 获取指定格式的当前时间字符串
func (l Ltime) NowTimeStr(format string) string {
	return time.Now().Format(format)
}

// NowTimePtr
// 获取当前时间指针
func (l Ltime) NowTimePtr() *time.Time {
	now := time.Now()
	return &now
}

// LocalTime 自定义时间类型，兼容 PostgreSQL
type LocalTime time.Time

// MarshalJSON LocalTime 实现 json 序列化接口
func (l LocalTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(l).Format(FormatDateTimeMs))
	return []byte(stamp), nil
}

// UnmarshalJSON LocalTime 实现 json 反序列化接口
func (l *LocalTime) UnmarshalJSON(data []byte) error {
	timeStr := Str{}.StringUnquote(strings.TrimSpace(string(data)))
	if timeStr = strings.TrimSpace(timeStr); timeStr == "" {
		return errors.New("时间字符串不能为空！")
	}

	jsonTime, err := time.Parse(GetTimeFormat(timeStr), timeStr)
	*l = LocalTime(jsonTime)
	return err
}

// Value LocalTime 写入数据库之前，转换成 time.Time
func (l LocalTime) Value() (driver.Value, error) {
	return time.Time(l).Format(FormatDateTimeMs), nil
}

// Scan 从数据库中读取数据，转换成 LocalTime
func (l *LocalTime) Scan(v interface{}) error {
	var sqlTime time.Time
	switch vt := v.(type) {
	case string:
		// 字符串转成 time.Time 类型
		sqlTime, _ = time.Parse(FormatDateTimeMs, vt)
	case time.Time:
		sqlTime = vt
	case *time.Time:
		sqlTime = *vt
	default:
		return errors.New("读取 LocalTime 类型处理错误！")
	}
	*l = LocalTime(sqlTime)
	return nil
}

// Sub 实现 time.Time.Sub 方法计算 LocalTime 时间差
func (l LocalTime) Sub(t LocalTime) time.Duration {
	return time.Time(l).Sub(time.Time(t))
}

// LocalDate 自定义日期类型，兼容 PostgreSQL
type LocalDate time.Time

// MarshalJSON LocalDate 实现 json 序列化接口
func (l LocalDate) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(l).Format(FormatDate))
	return []byte(stamp), nil
}

// UnmarshalJSON LocalDate 实现 json 反序列化接口
func (l *LocalDate) UnmarshalJSON(data []byte) error {
	dateStr := Str{}.StringUnquote(strings.TrimSpace(string(data)))
	if dateStr = strings.TrimSpace(dateStr); dateStr == "" {
		return errors.New("日期字符串不能为空！")
	}
	if len(dateStr) > len(FormatDate) && strings.Contains(dateStr, " ") {
		dateStr = strings.SplitN(dateStr, " ", 2)[0]
	}

	jsonTime, err := time.Parse(FormatDate, dateStr)
	*l = LocalDate(jsonTime)
	return err
}

// Value LocalDate 写入数据库之前，转换成 time.Time
func (l LocalDate) Value() (driver.Value, error) {
	if &l == nil {
		return nil, nil
	}
	return time.Time(l).Format(FormatDate), nil
}

// Scan 从数据库中读取数据，转换成 LocalDate
func (l *LocalDate) Scan(v interface{}) error {
	var sqlTime time.Time
	switch vt := v.(type) {
	case string:
		// 字符串转成 time.Time 类型
		sqlTime, _ = time.Parse(FormatDate, vt)
	case time.Time:
		sqlTime = vt
	case *time.Time:
		sqlTime = *vt
	default:
		return errors.New("读取 LocalDate 类型处理错误！")
	}
	*l = LocalDate(sqlTime)
	return nil
}

// Sub 实现 time.Time.Sub 方法计算 LocalDate 时间差
func (l LocalDate) Sub(t LocalDate) time.Duration {
	return time.Time(l).Sub(time.Time(t))
}

// NowLocalTime 获取当前自定义格式时间
func NowLocalTime() LocalTime {
	return LocalTime(time.Now())
}

// NowLocalTimePtr 获取当前自定义格式时间
func NowLocalTimePtr() *LocalTime {
	localTime := LocalTime(time.Now())
	return &localTime
}

// NowLocalDate 获取当前自定义格式日期
func NowLocalDate() LocalDate {
	return LocalDate(time.Now())
}

// 返回当前时间字符串
func NowlocalDatetimeStr() string {
	return GetTimeFormat(time.Now().Format(FormatDateTimeMs))
}

// NowLocalDatePtr NowLocalDate 获取当前自定义格式日期
func NowLocalDatePtr() *LocalDate {
	localDate := LocalDate(time.Now())
	return &localDate
}

// GetTimeFormat 根据日期时间字符串获取日期时间格式
func GetTimeFormat(timeStr string) (format string) {
	switch len(timeStr) {
	case 4:
		format = FormatYear
	case 5:
		format = FormatTime
	case 8:
		if strings.Contains(timeStr, ":") {
			format = FormatTimeSec
		} else if strings.Contains(timeStr, "-") {
			format = FormatDateShort
		} else {
			format = FormatDateNum
		}
	case 9:
		if strings.Contains(timeStr, "-") {
			format = FormatDateShort
		}
	case 10:
		if strings.Contains(timeStr, "-") &&
			!strings.Contains(timeStr, ":") &&
			!strings.Contains(timeStr, ".") {
			format = FormatDate
		} else {
			format = FormatTimeMs
		}
	case 11, 12:
		format = FormatTimeMs
	case 16:
		// yyyy-MM-dd HH:mm
		format = FormatDateTime
	case 19:
		// yyyy-MM-dd HH:mm:ss
		format = FormatDateTimeSec
	case 21, 22, 23:
		// yyyy-MM-dd HH:mm:ss.SSS
		format = FormatDateTimeMs
	default:
	}
	return
}

/**
 * @Description: 长时间
 */
type LongDateTime struct {
	time.Time
}

// UnmarshalJSON LongDateTime 实现 json 反序列化接口
func (ld *LongDateTime) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+FormatDateTimeMs+`"`, string(data), time.Local)
	*ld = LongDateTime{Time: now}
	return
}

// MarshalJSON LongDateTime 实现json 序列化
func (ld LongDateTime) MarshalJSON() ([]byte, error) {
	tTime := time.Time(ld.Time)
	if tTime.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%v\"", tTime.Format(FormatDateTimeMs))), nil
}

func (t LongDateTime) String() string {
	return time.Time(t.Time).Format(FormatDateTimeMs)
}

func (ld LongDateTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	tlt := time.Time(ld.Time)
	//判断给定时间是否和默认零时间的时间戳相同
	if tlt.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return tlt, nil
}

func (ld *LongDateTime) Scan(v interface{}) error {
	if value, ok := v.(time.Time); ok {
		*ld = LongDateTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

func (ld LongDateTime) GobEncode() ([]byte, error) {
	// 序列化 LongDateTime 的 UnixNano 值
	return gobEncodeInt64(ld.Time.Unix()), nil
}

func (ld *LongDateTime) GobDecode(data []byte) error {
	// 反序列化 UnixNano 值
	nano, err := gobDecodeInt64(data)
	if err != nil {
		return err
	}
	// 将 UnixNano 转换为 LongDateTime
	ld.Time = time.Unix(0, nano)
	return nil
}

// 辅助函数：将 int64 转换为 byte 切片
func gobEncodeInt64(n int64) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, n)
	return buf.Bytes()
}

// 辅助函数：将 byte 切片转换为 int64
func gobDecodeInt64(data []byte) (int64, error) {
	buf := bytes.NewBuffer(data)
	var n int64
	err := binary.Read(buf, binary.BigEndian, &n)
	return n, err
}
