คู่มือการใช้งาน pkg/presenter/jsend (หัวแปลงปลั๊ก)
แพ็กเกจ presenter คือที่อยู่ของ "หัวแปลงปลั๊ก" หรือเครื่องมือพิเศษที่เราสร้างขึ้นมาเพื่อ แปลง รูปแบบ Response มาตรฐานของเรา (Dashboard Style) ให้กลายเป็นรูปแบบอื่นที่ระบบภายนอก (เช่น Partner API) ต้องการ โดยไม่กระทบกับมาตรฐานหลักของโปรเจกต์

เอกสารนี้จะอธิบายการใช้งาน JSend Presenter โดยเฉพาะ

1. ปรัชญา: "กฎ" vs "ข้อยกเว้น"
pkg/response: คือ "กฎ" และมาตรฐานหลักของ API เรา ทุก Endpoint ปกติจะต้องใช้รูปแบบนี้

pkg/presenter: คือ "ข้อยกเว้น" ที่เราสร้างขึ้นมาสำหรับคุยกับคนอื่นที่เราควบคุมไม่ได้

เราจะ ไม่ เรียกใช้ Presenter ใน Handler ทั่วไป แต่จะใช้เฉพาะใน Handler ที่ถูกสร้างขึ้นมาสำหรับงานนั้นๆ โดยเฉพาะ

2. รูปแบบ JSend Standard
JSend คือมาตรฐานการตอบกลับ JSON ที่เรียบง่าย โดยใช้ status field เป็นตัวกำหนดสถานะของ Response

✅ JSend Success Response
{
    "status": "success",
    "data": {
        "id": 1,
        "name": "Nipon"
    }
}

status: จะเป็น "success" เสมอ

data: คือ "เนื้อ" ของข้อมูล

❌ JSend Error/Fail Response
{
    "status": "error",
    "message": "อีเมลนี้ถูกใช้งานแล้ว",
    "code": "ALREADY_EXISTS",
    "data": null
}

status: จะเป็น "error" (สำหรับ Server Error) หรือ "fail" (สำหรับ Client Error)

message: ข้อความที่มนุษย์อ่านเข้าใจได้

code & data: (Optional) ข้อมูลเพิ่มเติมสำหรับ Error

3. การใช้งาน "หัวแปลงปลั๊ก"
"หัวแปลงปลั๊ก" ของเราคือฟังก์ชัน Helper 2 ตัวที่อยู่ใน pkg/presenter/jsend.go

3.1 presenter.ToJSendSuccess(data)
ใช้สำหรับแปลงข้อมูลที่สำเร็จแล้ว ให้อยู่ในซอง JSend

// 1. เตรียมข้อมูล Response DTO ปกติของเราก่อน
responsePayload := toResponse(createdUserDomain)

// 2. เรียกใช้ "หัวแปลงปลั๊ก"
jsendResponse := presenter.ToJSendSuccess(responsePayload)

// 3. ส่งกลับไปเป็น JSON
return c.Status(fiber.StatusCreated).JSON(jsendResponse)

3.2 presenter.ToJSendError(err)
ใช้สำหรับแปลง *custom_errors.AppError มาตรฐานของเรา ให้อยู่ในซอง JSend

// `serviceErr` คือ *custom_errors.AppError ที่ได้มาจาก Service
if serviceErr != nil {
    // 1. เรียกใช้ "หัวแปลงปลั๊ก" สำหรับ Error
    jsendError := presenter.ToJSendError(serviceErr.(*custom_errors.AppError))
    
    // 2. ส่งกลับไปเป็น JSON
    return c.Status(serviceErr.HTTPStatus).JSON(jsendError)
}

4. ตัวอย่างการใช้งานจริงใน Handler พิเศษ
นี่คือตัวอย่างการสร้าง Endpoint /api/v1/partner/users ที่จะตอบกลับเป็น JSend format เสมอ

// ใน internal/modules/example_user/example_user_handler.go

import "go-template/pkg/presenter"

// CreateUserForPartner คือ Handler พิเศษสำหรับ Partner
func (h *handler) CreateUserForPartner(c fiber.Ctx) error {
	// --- ส่วนของ Logic เหมือนกับ Handler ปกติทุกประการ ---
	req := new(CreateRequest)
	if err := c.Bind().Body(req); err != nil {
		appErr := custom_errors.InvalidFormatError("Request body is not valid JSON", err.Error())
		// ⭐️ แปลง Error เป็น JSend ก่อนส่งกลับ ⭐️
		return c.Status(appErr.HTTPStatus).JSON(presenter.ToJSendError(appErr))
	}
	// ... (Validation Logic) ...

	domainData := &Domain{Name: req.Name, Email: req.Email}
	createdUserDomain, serviceErr := h.service.CreateUser(domainData, req.Password)
	if serviceErr != nil {
		appErr := serviceErr.(*custom_errors.AppError)
		// ⭐️ แปลง Error เป็น JSend ก่อนส่งกลับ ⭐️
		return c.Status(appErr.HTTPStatus).JSON(presenter.ToJSendError(appErr))
	}

	// --- ถึงเวลาใช้ "หัวแปลงปลั๊ก" ---
	responsePayload := toResponse(createdUserDomain)

	// ⭐️ แปลง Success Response ของเราให้เป็น JSend ⭐️
	jsendResponse := presenter.ToJSendSuccess(responsePayload)

	return c.Status(fiber.StatusCreated).JSON(jsendResponse)
}
