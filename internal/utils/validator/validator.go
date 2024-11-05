package validator

import (
	"github.com/asaskevich/govalidator"
)

func Check(someStruct interface{}) error {
	_, err := govalidator.ValidateStruct(someStruct)

	return err
}
