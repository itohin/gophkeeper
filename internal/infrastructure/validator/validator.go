package validator

import "net/mail"

func ValidateEmail() func(string) error {
	return func(email string) error {
		_, err := mail.ParseAddress(email)
		if err != nil {
			return err
		}
		return nil
	}
}
