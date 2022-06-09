package validator

import (
	"strings"

	"github.com/praveennagaraj97/online-consultation/interfaces"
)

func isPhoneNumberValida(phone interfaces.PhoneType) bool {

	var numberLen = len(strings.Join(strings.Split(phone.Number, ""), ""))

	switch strings.Replace(phone.Code, " ", "+", 1) {
	case "+91":
		return numberLen == 10
	default:
		return true
	}
}
