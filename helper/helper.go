package helper

import (
	"github.com/go-playground/validator/v10"
)

type Response struct {
	Meta meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

func APIResponse(message string, code int, status string, data interface{}) Response {
	meta := meta{
		Message: message,
		Code:    code,
		Status:  status,
	}
	jsonresponse := Response{
		Meta: meta,
		Data: data,
	}
	return jsonresponse
}

// func FormatValidationError(err error) []string {
// 	var errors []string

// 	for _, e := range err.(validator.ValidationErrors) {
// 		errors = append(errors, e.Error())
// 	}
// 	return errors
// }

// FormatValidationError mengubah error validasi menjadi pesan yang lebih informatif
func FormatValidationError(err error) []string {
	var errors []string
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, e := range ve {
			errors = append(errors, e.Error())
		}
	} else {
		errors = append(errors, err.Error())
	}
	return errors
}
