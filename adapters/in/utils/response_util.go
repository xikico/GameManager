package utils

import "time"

type Response[T any] struct {
	Status  int16  `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data"`
	TmpTime string `json:"tmp_time"` // 请求开始处理的时间
}

func Success[T any](data T, mes ...string) Response[T] {
	message := "操作成功"
	if len(mes) > 0 {
		message = mes[0]
	}
	startStr := time.Now().Format("2006-01-02 15:04:05")
	return Response[T]{
		Status:  200,
		Message: message,
		Data:    data,
		TmpTime: startStr,
	}
}

func Failure[T any](data T, mes ...string) Response[T] {
	message := "操作失败"
	if len(mes) > 0 {
		message = mes[0]
	}
	startStr := time.Now().Format("2006-01-02 15:04:05")
	return Response[T]{
		Status:  400,
		Message: message,
		Data:    data,
		TmpTime: startStr,
	}
}
