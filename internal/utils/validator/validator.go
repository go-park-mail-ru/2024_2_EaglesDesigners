package validator

import (
	"html"

	"github.com/asaskevich/govalidator"
)

func Check(someStruct interface{}) error {
	_, err := govalidator.ValidateStruct(someStruct)

	return err
}

func EscapePtrString(str *string) *string {
	var safe *string

	if str != nil {
		safe = new(string)
		*safe = html.EscapeString(*str)
	}

	return safe
}
