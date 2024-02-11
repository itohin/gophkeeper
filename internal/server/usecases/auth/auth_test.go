package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/itohin/gophkeeper/internal/server/entities"
	"github.com/itohin/gophkeeper/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuthUseCase_Register(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockUsersStorage(ctrl)
	hash := mocks.NewMockPasswordHasher(ctrl)
	otp := mocks.NewMockOTPGenerator(ctrl)
	mailer := mocks.NewMockMailer(ctrl)
	uuid := mocks.NewMockUUIDGenerator(ctrl)

	auth := &AuthUseCase{
		hash:   hash,
		uuid:   uuid,
		otp:    otp,
		repo:   repo,
		mailer: mailer,
	}

	email := "email@mall.ru"
	password := "password"
	passwordHash := "password_hash"
	var id [16]byte
	copy(id[:], "1955a7d6-0968-425b-bdb6-fb9a0e4b39e7")
	otpCode := "otp_secret"
	ctx := context.TODO()
	user := entities.User{
		ID:               id,
		Email:            email,
		Password:         passwordHash,
		VerificationCode: otpCode,
	}

	tests := []struct {
		name      string
		mockTimes map[string]int
		errors    map[string]error
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "PasswordHash error",
			mockTimes: map[string]int{
				"hash":   1,
				"uuid":   0,
				"otp":    0,
				"repo":   0,
				"mailer": 0,
			},
			errors: map[string]error{
				"hash":   errors.New("password hash error"),
				"uuid":   nil,
				"otp":    nil,
				"repo":   nil,
				"mailer": nil,
			},
			wantErr: assert.Error,
		},
		{
			name: "UUID error",
			mockTimes: map[string]int{
				"hash":   1,
				"uuid":   1,
				"otp":    0,
				"repo":   0,
				"mailer": 0,
			},
			errors: map[string]error{
				"hash":   nil,
				"uuid":   errors.New("uuid error"),
				"otp":    nil,
				"repo":   nil,
				"mailer": nil,
			},
			wantErr: assert.Error,
		},
		{
			name: "OTP error",
			mockTimes: map[string]int{
				"hash":   1,
				"uuid":   1,
				"otp":    1,
				"repo":   0,
				"mailer": 0,
			},
			errors: map[string]error{
				"hash":   nil,
				"uuid":   nil,
				"otp":    errors.New("otp error"),
				"repo":   nil,
				"mailer": nil,
			},
			wantErr: assert.Error,
		},
		{
			name: "OTP error",
			mockTimes: map[string]int{
				"hash":   1,
				"uuid":   1,
				"otp":    1,
				"repo":   1,
				"mailer": 0,
			},
			errors: map[string]error{
				"hash":   nil,
				"uuid":   nil,
				"otp":    nil,
				"repo":   errors.New("repo error"),
				"mailer": nil,
			},
			wantErr: assert.Error,
		},
		{
			name: "Mailer error",
			mockTimes: map[string]int{
				"hash":   1,
				"uuid":   1,
				"otp":    1,
				"repo":   1,
				"mailer": 1,
			},
			errors: map[string]error{
				"hash":   nil,
				"uuid":   nil,
				"otp":    nil,
				"repo":   nil,
				"mailer": errors.New("mailer error"),
			},
			wantErr: assert.Error,
		},
		{
			name: "Success",
			mockTimes: map[string]int{
				"hash":   1,
				"uuid":   1,
				"otp":    1,
				"repo":   1,
				"mailer": 1,
			},
			errors: map[string]error{
				"hash":   nil,
				"uuid":   nil,
				"otp":    nil,
				"repo":   nil,
				"mailer": nil,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			hash.EXPECT().HashPassword(password).Return(passwordHash, tt.errors["hash"]).Times(tt.mockTimes["hash"])
			uuid.EXPECT().Generate().Return(id, tt.errors["uuid"]).Times(tt.mockTimes["uuid"])
			otp.EXPECT().RandomSecret().Return(otpCode, tt.errors["otp"]).Times(tt.mockTimes["otp"])
			repo.EXPECT().Save(ctx, user).Return(user, tt.errors["repo"]).Times(tt.mockTimes["repo"])
			mailer.EXPECT().SendMail([]string{email}, "Для подтверждения адреса электронной почты в сервисе gophkeeper, введите пожалуйста код подтверждения: "+otpCode).Return(tt.errors["mailer"]).Times(tt.mockTimes["mailer"])

			tt.wantErr(t, auth.Register(ctx, email, password), fmt.Sprintf("Register()"))
		})
	}
}
