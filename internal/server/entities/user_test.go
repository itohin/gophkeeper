package entities

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUser_Verify_Success(t *testing.T) {

	u := &User{
		VerificationCode: "secret_otp",
	}

	err := u.Verify("secret_otp")
	if err != nil {
		t.Errorf("Verify() success error = %v", err)
	}
	assert.Equal(t, true, u.VerifiedAt.Valid)
}

func TestUser_Verify_Errors(t *testing.T) {

	type args struct {
		otp            string
		isUserVerified bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr string
	}{
		{
			name:    "already verified",
			args:    args{isUserVerified: true, otp: "otp_secret"},
			wantErr: "already verified",
		},
		{
			name:    "wrong otp",
			args:    args{isUserVerified: false, otp: "wrong_otp_secret"},
			wantErr: "wrong verification code",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				VerificationCode: "secret_otp",
				VerifiedAt: sql.NullTime{
					Time:  time.Now(),
					Valid: tt.args.isUserVerified,
				},
			}

			err := u.Verify(tt.args.otp)
			if err.Error() != tt.wantErr {
				t.Errorf("Verify() errors error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUser_IsVerified(t *testing.T) {
	tests := []struct {
		name       string
		VerifiedAt sql.NullTime
		want       bool
	}{
		{
			name: "verified",
			VerifiedAt: sql.NullTime{
				Valid: true,
			},
			want: true,
		},
		{
			name: "not verified",
			VerifiedAt: sql.NullTime{
				Valid: false,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				VerifiedAt: tt.VerifiedAt,
			}
			assert.Equalf(t, tt.want, u.IsVerified(), "IsVerified()")
		})
	}
}
