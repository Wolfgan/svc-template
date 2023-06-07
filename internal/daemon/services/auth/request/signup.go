package request

import (
	"service-template/internal/config/valid"
	"service-template/internal/model"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// SignUp Структура HTTP-запроса на регистрацию пользователя
type SignUp struct {
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Password string `json:"password"`
}

func (in SignUp) Validate() error {
	return validation.ValidateStruct(&in,
		//validation.Field(&in.Email, is.Email),
		//validation.Field(&in.Phone, is.Digit),
		validation.Field(&in.Phone, validation.Required.When(in.Email == "").Error("either phone or email is required.")),
		validation.Field(&in.Email, validation.Required.When(in.Phone == "").Error("either phone or email is required.")),
		validation.Field(&in.Password, validation.Required, validation.Match(valid.Password)),
	)
}

func (in SignUp) ToModel() *model.User {
	return &model.User{
		Email:    in.Email,
		Phone:    in.Phone,
		Password: in.Password,
	}
}
