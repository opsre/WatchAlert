package tools

import (
	"context"
	"fmt"
	"reflect"
	"strconv"

	"github.com/mitchellh/mapstructure"
	"github.com/zeromicro/go-zero/core/logc"
)

func ConvertStringToInt(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		logc.Error(context.Background(), fmt.Sprintf("Convert String to int failed, err: %s", err.Error()))
		return 0
	}

	return num
}

func ConvertStringToInt64(str string) int64 {
	num64, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		logc.Error(context.Background(), fmt.Sprintf("Convert String to int64 failed, err: %s", err.Error()))
		return 0
	}

	return num64
}

func ConvertStructToMap(v interface{}) map[string]interface{} {
	var result map[string]interface{} // 可为 nil，解码器会分配

	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:           &result, // 指向 result
		TagName:          "json",  // 读取 json 标签
		WeaklyTypedInput: true,    // 需要宽松转换时打开
		// Squash:        true,    // 有匿名嵌入字段时可打开
		// ZeroFields:    true,    // 需要用零值覆盖时可打开
	})
	if err != nil {
		logc.Error(context.Background(), "ConvertStructToMap NewDecoder failed", "err", err)
		return nil
	}

	if err := dec.Decode(v); err != nil {
		logc.Error(context.Background(), "ConvertStructToMap Decode failed", "err", err)
		return nil
	}
	return result
}

func ConvertSliceToMapList(slice interface{}) []map[string]interface{} {
	// 1. 检查输入是否为切片
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		logc.Error(context.Background(), "ConvertSliceToMapList input is not a slice")
		return nil
	}

	length := v.Len()
	if length == 0 {
		return []map[string]interface{}{}
	}

	// 2. 预分配结果切片
	result := make([]map[string]interface{}, 0, length)

	// 3. 遍历每个元素，转换为 map
	for i := 0; i < length; i++ {
		item := v.Index(i).Interface()
		m := ConvertStructToMap(item)
		if m == nil {
			// 可选：跳过失败项，或 append 空 map
			m = make(map[string]interface{})
		}
		result = append(result, m)
	}

	return result
}
