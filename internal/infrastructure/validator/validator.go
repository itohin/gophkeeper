package validator

import (
	"fmt"
	"net/mail"
)

func ValidateEmail() func(string) error {
	return func(email string) error {
		_, err := mail.ParseAddress(email)
		if err != nil {
			return err
		}
		return nil
	}
}

func ValidateConfirmationCode(code string) func(string) error {
	return func(input string) error {
		if input != code {
			return fmt.Errorf("неправильный код подтверждения: %s/n", input)
		}
		return nil
	}
}
