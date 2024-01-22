package objectvaule

import "net/http"

type ResponseValue struct {
	Title   string
	Message string
	IsOk    bool
	Status  int32
	Value   any
}

func BadResponseSingle(message string) *ResponseValue {
	return &ResponseValue{
		Title:   "Proceso no existoso",
		IsOk:    false,
		Message: message,
		Status:  http.StatusBadRequest,
	}
}
