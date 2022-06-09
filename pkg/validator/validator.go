package validator

import (
	"net/mail"

	"github.com/praveennagaraj97/online-consultation/interfaces"
)

func ValidateEmail(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return err
	}

	return nil
}

func ValidatePhoneNumber(phone interfaces.PhoneType) *map[string]string {
	errs := map[string]string{}

	if phone.Code == "" {
		errs["code"] = "Phone Code cannot be empty"
	}

	if phone.Number == "" {
		errs["number"] = "Phone Number cannot be empty"
	}

	if isValid := isPhoneNumberValida(phone); !isValid {
		errs["phone_number"] = "Entered Phone Number is not valid"
	}

	if len(errs) > 0 {
		return &errs
	}

	return nil

}
