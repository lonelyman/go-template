คู่มือการใช้งาน pkg/validator
แพ็กเกจ validator นี้เป็นเครื่องมือตรวจสอบข้อมูล (Validation) ที่ทรงพลัง ถูกสร้างขึ้นมาโดยครอบ (wrap) ไลบรารี go-playground/validator/v10 ที่เป็นมาตรฐานอุตสาหกรรมเอาไว้

เป้าหมายหลัก: เพื่อให้การตรวจสอบข้อมูลใน DTO (Data Transfer Objects) ของเราง่าย, เป็นมาตรฐาน, และสามารถคืนค่า Error ที่มีโครงสร้างชัดเจนพร้อม ข้อความที่กำหนดเอง (Custom Error Message) ได้อย่างยืดหยุ่น

🚀 การใช้งานพื้นฐาน
หัวใจของแพ็กเกจนี้คือฟังก์ชัน validator.Validate() ซึ่งรับ struct ใดๆ เข้าไปตรวจสอบ และจะคืนค่าเป็น \*ValidationResult กลับมาเสมอ

1. การประกาศ DTO Struct
   ใน handler ของเรา ให้ประกาศ struct สำหรับรับข้อมูล (Request DTO) พร้อมกับ validate tag มาตรฐาน และ json tag สำหรับการ binding

// ในไฟล์ handler

type CreateUserRequest struct {
Name string `json:"name" validate:"required,min=2"`
Email string `json:"email" validate:"required,email"`
Password string `json:"password" validate:"required,min=8"`
}

2. การเรียกใช้งานใน Handler
   ในเมธอดของ handler เราจะเรียกใช้ validator.Validate() หลังจากที่ทำการ Bind() ข้อมูลเรียบร้อยแล้ว

import (
"go-template/pkg/validator"
"go-template/pkg/custom_errors"
"go-template/pkg/response"
)

func (h \*handler) CreateUser(c fiber.Ctx) error {
// 1. Bind ข้อมูลจาก Request Body เข้า DTO
req := new(CreateRequest)
if err := c.Bind().Body(req); err != nil {
// ... จัดการ Binding Error ...
}

    // 2. เรียกใช้ Validator ของเรา
    if validationResult := validator.Validate(req); !validationResult.IsValid {
        // 3. ถ้าไม่ผ่าน ให้สร้าง ValidationError พร้อมรายละเอียด
        appErr := custom_errors.ValidationError("ข้อมูลที่ส่งมาไม่ถูกต้อง", validationResult.Errors)
        return response.Error(c, appErr)
    }

    // ... ถ้าผ่าน ก็ทำงานต่อไป ...
    return response.Success(c, fiber.StatusCreated, "User created", nil)

}

✨ การใช้งานขั้นสูง: Custom Error Messages ด้วย vmsg Tag
นี่คือความสามารถที่ทรงพลังที่สุดของ Validator ของเรา! เราสามารถกำหนดข้อความ Error สำหรับแต่ละกฎ (rule) ได้โดยตรงใน Struct Tag โดยใช้ vmsg

รูปแบบ (Syntax)
vmsg:"rule1:message1,rule2:message2,..."

ใช้ , (comma) ในการคั่นระหว่างแต่ละกฎ

ใช้ : (colon) ในการคั่นระหว่างชื่อกฎกับข้อความ

ถ้าต้องการใช้ , ภายในข้อความ ให้ทำการ escape ด้วย \ (เช่น \,)

สามารถใช้ " และ ' ในข้อความได้ตามปกติ

ตัวอย่างการใช้งานเต็มรูปแบบ
type RegisterRequest struct {
// ตัวอย่างที่ 1: กำหนดข้อความสำหรับแต่ละ Rule แยกกัน
FullName string `json:"fullName" validate:"required,min=2"
	                 vmsg:"required:กรุณากรอกชื่อ-นามสกุล,min:ชื่อต้องยาวอย่างน้อย 2 ตัวอักษร"`

    // ตัวอย่างที่ 2: ใช้ Single Quote (') ในข้อความ
    Username string `json:"username" validate:"required,alphanum"
                     vmsg:"required:กรุณากรอกชื่อผู้ใช้ (เช่น 'nipon.k')"`

    // ตัวอย่างที่ 3: ใช้ Double Quote (") ในข้อความ
    UserBio string `json:"bio" validate:"required"
                    vmsg:"required:กรุณาใส่คำอธิบายสั้นๆ, เช่น \"I love Go!\""`

    // ตัวอย่างที่ 4: ใช้คอมม่า (,) ในข้อความโดยการ escape ด้วย `\,`
    Address string `json:"address" validate:"required"
                     vmsg:"required:กรุณากรอกที่อยู่ (บ้านเลขที่\, ถนน\, ตำบล)"`

    // ตัวอย่างที่ 5: ใช้ Default Message (ไม่มี vmsg)
    Age int `json:"age" validate:"required,numeric"`

}

ผลลัพธ์ของ Error Response
ถ้าการ Validate ล้มเหลว validationResult.Errors จะถูกส่งไปให้ custom_errors.ValidationError และสุดท้ายจะถูกจัดรูปแบบโดย response.Error ทำให้ได้ JSON Response ที่สวยงามและเป็นประโยชน์ต่อ Frontend แบบนี้:

{
"success": false,
"message": "ข้อมูลที่ส่งมาไม่ถูกต้อง",
"error": {
"code": "VALIDATION_ERROR",
"details": [
{
"field": "fullName",
"message": "ชื่อต้องยาวอย่างน้อย 2 ตัวอักษร",
"value": "N"
},
{
"field": "age",
"message": "ฟิลด์นี้จำเป็นต้องระบุ",
"value": "0"
}
]
}
}

📚 พจนานุกรม Default Messages
ถ้าเราไม่กำหนด vmsg tag, Validator ของเราจะใช้ข้อความเริ่มต้น (Default) ที่มีความหมายดีเหล่านี้แทน:

Tag

Default Message

required

ฟิลด์นี้จำเป็นต้องระบุ

email

ต้องเป็นรูปแบบอีเมลที่ถูกต้อง

url

ต้องเป็น URL ที่ถูกต้อง

uuid

ต้องเป็น UUID ที่ถูกต้อง

min

ต้องมีขนาดอย่างน้อย {param}

max

ต้องมีขนาดไม่เกิน {param}

len

ต้องมีขนาดเท่ากับ {param} พอดี

numeric

ต้องเป็นตัวเลขเท่านั้น

gt

ต้องมีค่ามากกว่า {param}

gte

ต้องมีค่าอย่างน้อย {param}

lt

ต้องมีค่าน้อยกว่า {param}

lte

ต้องมีค่าไม่เกิน {param}

alphanum

ต้องเป็นตัวอักษรหรือตัวเลขเท่านั้น

alpha

ต้องเป็นตัวอักษรเท่านั้น

datetime

ต้องเป็นวันที่และเวลาในรูปแบบที่ถูกต้อง ({param})

(อื่นๆ)

จะแสดง Error ดั้งเดิมจากไลบรารี
