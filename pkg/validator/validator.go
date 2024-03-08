package validator

import (
	"errors"
	"fmt"
	"net/mail"
	"unicode"
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

func ValidateConfirmationCode() func(string) error {
	return func(input string) error {
		if len(input) < 1 {
			return fmt.Errorf("неправильный код подтверждения: %s/n", input)
		}
		return nil
	}
}

func ValidatePassword() func(string) error {
	return func(password string) error {
		if !isValidPassword(password) {
			return errors.New("пароль должен быть не короче 8 символов, строчные и прописные буквы, цифры, спецсимволы")
		}
		return nil
	}
}

func isValidPassword(password string) bool {
	if len(password) < 8 || len(password) > 127 {
		return false
	}
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasUpper && hasLower && hasNumber && hasSpecial
}
