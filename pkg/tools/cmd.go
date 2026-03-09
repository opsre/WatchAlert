package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"regexp"
	"time"

	"github.com/bytedance/sonic"
	"github.com/rs/xid"
	"github.com/zeromicro/go-zero/core/logc"
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

// IsJSON 判断字符串是否为 JSON 格式
func IsJSON(str string) bool {
	var js map[string]interface{}
	return sonic.Unmarshal([]byte(str), &js) == nil
}

// FormatJson 格式化 JSON 字符串
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
