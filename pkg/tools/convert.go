package tools

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/zeromicro/go-zero/core/logc"
	"strconv"
	"watchAlert/internal/models"
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

func ConvertEventToMap(event models.AlertCurEvent) map[string]interface{} {
	data := make(map[string]interface{})
	eventJson := JsonMarshalToByte(event)
	err := sonic.Unmarshal(eventJson, &data)
	if err != nil {
		logc.Error(context.Background(), "ConvertEventToMap Unmarshal failed: ", err)
	}

	return data
}
