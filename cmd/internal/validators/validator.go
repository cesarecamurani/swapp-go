package validators

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/nyaruka/phonenumbers"
)

func Init() {
	if val, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = val.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
			phone := fl.Field().String()

			num, err := phonenumbers.Parse(phone, "")
			if err != nil {
				return false
			}

			return phonenumbers.IsValidNumber(num)
		})
	}
}
