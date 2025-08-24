// pkg/presenter/jsend.go
package presenter

import "go-template/pkg/custom_errors"

// JSendSuccess คือพิมพ์เขียวสำหรับ Success Response ในรูปแบบ JSend
type JSendSuccess struct {
	Status string      `json:"status"` // จะมีค่าเป็น "success" เสมอ
	Data   interface{} `json:"data"`
}

// JSendError คือพิมพ์เขียวสำหรับ Error Response ในรูปแบบ JSend
type JSendError struct {
	Status  string      `json:"status"` // จะมีค่าเป็น "error" เสมอ
	Message string      `json:"message"`
	Code    string      `json:"code,omitempty"`
	Data    interface{} `json:"data,omitempty"` // สำหรับ validation details
}

// ToJSendSuccess คือ "เครื่องมือ" ที่ใช้แปลงข้อมูลใดๆ ให้อยู่ในรูปแบบ JSend Success
// มันจะรับข้อมูล (data) เข้ามา แล้วห่อด้วยซอง JSend ให้
func ToJSendSuccess(data interface{}) JSendSuccess {
	return JSendSuccess{
		Status: "success",
		Data:   data,
	}
}

// ToJSendError คือ "เครื่องมือ" ที่ใช้แปลง AppError มาตรฐานของเรา ให้อยู่ในรูปแบบ JSend Error
func ToJSendError(err *custom_errors.AppError) JSendError {
	return JSendError{
		Status:  "error",
		Message: err.Message,
		Code:    err.Code,
		Data:    err.Details,
	}
}
