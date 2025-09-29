package tools

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/rs/xid"
	"github.com/zeromicro/go-zero/core/logc"
	"io"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func RandId() string {
	return xid.New().String()
}

func RandUid() string {
	limit := 8
	gid := xid.New().String()

	var xx []string
	for _, v := range gid {
		xx = append(xx, string(v))
	}

	rand.Seed(time.Now().UnixNano())
	perm := rand.Perm(len(gid))

	var id string
	for i := 0; i < limit; i++ {
		id += xx[perm[i]]
	}

	return id
}

func JsonMarshalToString(v interface{}) string {
	data, err := sonic.Marshal(v)
	if err != nil {
		logc.Error(context.Background(), fmt.Sprintf("序列化失败, 原数据「%v」, err: %s", v, err.Error()))
		return "{}"
	}
	return string(data)
}

func JsonMarshalToByte(v interface{}) []byte {
	data, err := sonic.Marshal(v)
	if err != nil {
		logc.Error(context.Background(), fmt.Sprintf("序列化失败, 原数据「%v」, err: %s", v, err.Error()))
		return []byte("{}")
	}
	return data
}

// ParserVariables 处理告警内容中变量形式的字符串，替换为对应的值
func ParserVariables(annotations string, data map[string]interface{}) string {
	// 使用正则表达式匹配变量形式的字符串
	re := regexp.MustCompile(`\$\{(.*?)\}`)

	// 使用正则表达式替换所有匹配的变量
	return re.ReplaceAllStringFunc(annotations, func(match string) string {
		variable := match[2 : len(match)-1] // 提取变量名
		return fmt.Sprintf("%v", getJSONValue(data, variable))
	})
}

// 通过变量形式 ${key} 获取 JSON 数据中的值
func getJSONValue(data map[string]interface{}, variable string) interface{} {
	// 使用正则表达式分割键名数组
	keys := regexp.MustCompile(`\.`).Split(variable, -1)

	// 逐级获取 JSON 数据中的值
	for _, key := range keys {
		if v, ok := data[key]; ok {
			data, ok = v.(map[string]interface{})
			if !ok {
				return v
			}
		} else {
			return nil
		}
	}

	return nil
}

func IsJSON(str string) bool {
	var js map[string]interface{}
	return sonic.Unmarshal([]byte(str), &js) == nil
}

func FormatJson(s string) string {
	var ns string
	if IsJSON(s) {
		// 将字符串解析为map类型
		var data map[string]interface{}
		err := sonic.Unmarshal([]byte(s), &data)
		if err != nil {
			logc.Errorf(context.Background(), fmt.Sprintf("Error parsing JSON: %s", err.Error()))
		} else {
			// 格式化JSON并输出
			formattedJson, err := json.MarshalIndent(data, "", "  ")
			if err != nil {
				logc.Errorf(context.Background(), fmt.Sprintf("Error marshalling JSON: %s", err.Error()))
			} else {
				ns = string(formattedJson)
			}
		}
	} else {
		// 不是 json 格式的需要转义下其中的特殊符号，并且只取双引号(")内的内容。
		ns = JsonMarshalToString(s)
		ns = ns[1 : len(ns)-1]
	}
	return ns
}

// ParseReaderBody 处理请求Body
func ParseReaderBody(body io.Reader, req interface{}) error {
	newBody := body
	bodyByte, err := io.ReadAll(newBody)
	if err != nil {
		return fmt.Errorf("读取 Body 失败, err: %s", err.Error())
	}
	if err := sonic.Unmarshal(bodyByte, &req); err != nil {
		return fmt.Errorf("解析 Body 失败, body: %s, err: %s", string(bodyByte), err.Error())
	}
	return nil
}

func ParseTime(month string) (int, time.Month, int) {
	parsedTime, err := time.Parse("2006-01", month)
	if err != nil {
		return 0, time.Month(0), 0
	}
	curYear, curMonth, curDay := parsedTime.Date()
	return curYear, curMonth, curDay
}

func GetWeekday(date string) (time.Weekday, error) {
	t, err := time.Parse("2006-1-2", date)
	if err != nil {
		return 0, err
	}

	weekday := t.Weekday()
	return weekday, nil
}

func IsEndOfWeek(dateStr string) bool {
	date, err := time.Parse("2006-1-2", dateStr)
	if err != nil {
		return false
	}
	return date.Weekday() == time.Sunday
}

func ProcessRuleExpr(ruleExpr string) (operator string, value float64, err error) {
	var supportedOperators = []string{">=", "<=", "==", "!=", ">", "<", "="}

	// 去除表达式两端的空白字符
	trimmedExpr := strings.TrimSpace(ruleExpr)

	// 遍历操作符列表。
	for _, op := range supportedOperators {
		if strings.HasPrefix(trimmedExpr, op) {
			// 提取数值
			valueStr := strings.TrimPrefix(trimmedExpr, op)
			value, err = strconv.ParseFloat(strings.TrimSpace(valueStr), 64)
			if err != nil {
				return "", 0, fmt.Errorf("无法解析数值 '%s': %w", valueStr, err)
			}

			return op, value, nil
		}
	}

	return "", 0, fmt.Errorf("无效的表达式，未找到有效的操作符: %s", ruleExpr)
}

func GenerateHashPassword(passwd string) string {
	arr := md5.Sum([]byte(passwd))
	return hex.EncodeToString(arr[:])
}
