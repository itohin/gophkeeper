package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/itohin/gophkeeper/internal/server/entities"
	"github.com/itohin/gophkeeper/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAuthUseCase_Register(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usersRepo := mocks.NewMockUsersStorage(ctrl)
	hash := mocks.NewMockPasswordHasher(ctrl)
	otp := mocks.NewMockOTPGenerator(ctrl)
	mailer := mocks.NewMockMailer(ctrl)
	uuid := mocks.NewMockUUIDGenerator(ctrl)

	auth := &AuthUseCase{
		hash:      hash,
		uuid:      uuid,
		otp:       otp,
		usersRepo: usersRepo,
		mailer:    mailer,
	}

	email := "email@mall.ru"
	password := "password"
	passwordHash := "password_hash"
	var id [16]byte
	copy(id[:], "1955a7d6-0968-425b-bdb6-fb9a0e4b39e7")
	otpCode := "otp_secret"
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
				"hash": errors.New("password hash error"),
				"uuid": nil,
				"otp":  nil,
				"repo": nil,
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
				"hash": nil,
				"uuid": errors.New("uuid error"),
				"otp":  nil,
				"repo": nil,
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
				"hash": nil,
				"uuid": nil,
				"otp":  errors.New("otp error"),
				"repo": nil,
			},
			wantErr: assert.Error,
		},
		{
			name: "repo error",
			mockTimes: map[string]int{
				"hash":   1,
				"uuid":   1,
				"otp":    1,
				"repo":   1,
				"mailer": 0,
			},
			errors: map[string]error{
				"hash": nil,
				"uuid": nil,
				"otp":  nil,
				"repo": errors.New("repo error"),
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
				"hash": nil,
				"uuid": nil,
				"otp":  nil,
				"repo": nil,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			hash.EXPECT().HashPassword(password).Return(passwordHash, tt.errors["hash"]).Times(tt.mockTimes["hash"])
			uuid.EXPECT().Generate().Return(id, tt.errors["uuid"]).Times(tt.mockTimes["uuid"])
			otp.EXPECT().RandomSecret().Return(otpCode, tt.errors["otp"]).Times(tt.mockTimes["otp"])
			usersRepo.EXPECT().Save(gomock.Any(), user).Return(tt.errors["repo"]).Times(tt.mockTimes["repo"])
			mailer.EXPECT().SendMailAsync([]string{email}, "Для подтверждения адреса электронной почты в сервисе gophkeeper, введите пожалуйста код подтверждения: "+otpCode).Times(tt.mockTimes["mailer"])

			tt.wantErr(t, auth.Register(context.TODO(), email, password), fmt.Sprintf("Register()"))
		})
	}
}

func TestAuthUseCase_Verify(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usersRepo := mocks.NewMockUsersStorage(ctrl)
	uuid := mocks.NewMockUUIDGenerator(ctrl)
	jwtManager := mocks.NewMockJWTManager(ctrl)
	dbTransaction := mocks.NewMockDBTransactionManager(ctrl)

	auth := &AuthUseCase{
		uuid:      uuid,
		usersRepo: usersRepo,
		jwt:       jwtManager,
		tx:        dbTransaction,
	}

	var userId [16]byte
	var sessionId [16]byte
	copy(userId[:], "1955a7d6-0968-425b-bdb6-fb9a0e4b39e7")
	copy(sessionId[:], "1843a7d7-1268-345b-bdb9-ga3a0e4b34e8")
	refreshExpiration := time.Now().Add(time.Second * 10)

	tests := []struct {
		name      string
		mockTimes map[string]int
		errors    map[string]error
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "find user error",
			mockTimes: map[string]int{
				"users_repo_find": 1,
				"jwt":             0,
				"uuid":            0,
				"get_refresh_ttl": 0,
				"transaction":     0,
			},
			errors: map[string]error{
				"users_repo_find": errors.New("find user error"),
				"jwt":             nil,
				"uuid":            nil,
				"transaction":     nil,
			},
			wantErr: assert.Error,
		},
		{
			name: "make jwt error",
			mockTimes: map[string]int{
				"users_repo_find": 1,
				"jwt":             1,
				"uuid":            0,
				"get_refresh_ttl": 0,
				"transaction":     0,
			},
			errors: map[string]error{
				"users_repo_find": nil,
				"jwt":             errors.New("jwt error"),
				"uuid":            nil,
				"transaction":     nil,
			},
			wantErr: assert.Error,
		},
		{
			name: "make session uuid error",
			mockTimes: map[string]int{
				"users_repo_find": 1,
				"jwt":             1,
				"uuid":            1,
				"get_refresh_ttl": 0,
				"transaction":     0,
			},
			errors: map[string]error{
				"users_repo_find": nil,
				"jwt":             nil,
				"uuid":            errors.New("uuid error"),
				"transaction":     nil,
			},
			wantErr: assert.Error,
		},
		{
			name: "db transaction error",
			mockTimes: map[string]int{
				"users_repo_find": 1,
				"jwt":             1,
				"uuid":            1,
				"get_refresh_ttl": 1,
				"transaction":     1,
			},
			errors: map[string]error{
				"users_repo_find": nil,
				"jwt":             nil,
				"uuid":            nil,
				"transaction":     errors.New("db transaction error"),
			},
			wantErr: assert.Error,
		},
		{
			name: "success",
			mockTimes: map[string]int{
				"users_repo_find": 1,
				"jwt":             1,
				"uuid":            1,
				"get_refresh_ttl": 1,
				"transaction":     1,
			},
			errors: map[string]error{
				"users_repo_find": nil,
				"jwt":             nil,
				"uuid":            nil,
				"transaction":     nil,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			user := entities.User{
				ID:               userId,
				Email:            "email@mall.ru",
				VerificationCode: "otp_secret",
			}
			session := entities.Session{
				ID: sessionId,
			}

			usersRepo.EXPECT().FindByEmail(gomock.Any(), "email@mall.ru").Return(&user, tt.errors["users_repo_find"]).Times(tt.mockTimes["users_repo_find"])
			jwtManager.EXPECT().MakeJWT(user.ID.String()).Return("jwt.token", tt.errors["jwt"]).Times(tt.mockTimes["jwt"])
			uuid.EXPECT().Generate().Return(sessionId, tt.errors["uuid"]).Times(tt.mockTimes["uuid"])
			jwtManager.EXPECT().MakeRefreshExpiration().Return(refreshExpiration).Times(tt.mockTimes["get_refresh_ttl"])
			dbTransaction.EXPECT().Transaction(gomock.Any(), gomock.Any()).Return(tt.errors["transaction"]).Times(tt.mockTimes["transaction"])

			token, err := auth.Verify(context.TODO(), "email@mall.ru", "otp_secret", "unique_fingerprint")

			tt.wantErr(t, err, fmt.Sprintf("Verify()"))

			if token != nil {
				assert.Equal(t, token.RefreshToken, session.ID.String())
				assert.Equal(t, token.AccessToken, "jwt.token")
			}
		})
	}
}

func TestAuthUseCase_Login(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usersRepo := mocks.NewMockUsersStorage(ctrl)
	uuid := mocks.NewMockUUIDGenerator(ctrl)
	jwtManager := mocks.NewMockJWTManager(ctrl)
	hash := mocks.NewMockPasswordHasher(ctrl)
	sessionsRepo := mocks.NewMockSessionsStorage(ctrl)

	auth := &AuthUseCase{
		uuid:         uuid,
		usersRepo:    usersRepo,
		jwt:          jwtManager,
		sessionsRepo: sessionsRepo,
		hash:         hash,
	}

	var userId [16]byte
	var sessionId [16]byte
	copy(userId[:], "1955a7d6-0968-425b-bdb6-fb9a0e4b39e7")
	copy(sessionId[:], "1843a7d7-1268-345b-bdb9-ga3a0e4b34e8")
	refreshExpiration := time.Now().Add(time.Second * 10)

	tests := []struct {
		name            string
		mockTimes       map[string]int
		errors          map[string]error
		isPasswordValid bool
		wantErr         assert.ErrorAssertionFunc
	}{
		{
			name: "find user error",
			mockTimes: map[string]int{
				"users_repo_find":      1,
				"password_hash":        0,
				"jwt":                  0,
				"uuid":                 0,
				"get_refresh_ttl":      0,
				"sessions_repo_delete": 0,
				"sessions_repo_save":   0,
			},
			errors: map[string]error{
				"users_repo_find":      errors.New("find user error"),
				"jwt":                  nil,
				"uuid":                 nil,
				"sessions_repo_delete": nil,
				"sessions_repo_save":   nil,
			},
			isPasswordValid: true,
			wantErr:         assert.Error,
		},
		{
			name: "wrong password",
			mockTimes: map[string]int{
				"users_repo_find":      1,
				"password_hash":        1,
				"jwt":                  0,
				"uuid":                 0,
				"get_refresh_ttl":      0,
				"sessions_repo_delete": 0,
				"sessions_repo_save":   0,
			},
			errors: map[string]error{
				"users_repo_find":      nil,
				"jwt":                  nil,
				"uuid":                 nil,
				"sessions_repo_delete": nil,
				"sessions_repo_save":   nil,
			},
			isPasswordValid: false,
			wantErr:         assert.Error,
		},
		{
			name: "jwt error",
			mockTimes: map[string]int{
				"users_repo_find":      1,
				"password_hash":        1,
				"jwt":                  1,
				"uuid":                 0,
				"get_refresh_ttl":      0,
				"sessions_repo_delete": 0,
				"sessions_repo_save":   0,
			},
			errors: map[string]error{
				"users_repo_find":      nil,
				"jwt":                  errors.New("make jwt error"),
				"uuid":                 nil,
				"sessions_repo_delete": nil,
				"sessions_repo_save":   nil,
			},
			isPasswordValid: true,
			wantErr:         assert.Error,
		},
		{
			name: "session uuid error",
			mockTimes: map[string]int{
				"users_repo_find":      1,
				"password_hash":        1,
				"jwt":                  1,
				"uuid":                 1,
				"get_refresh_ttl":      0,
				"sessions_repo_delete": 0,
				"sessions_repo_save":   0,
			},
			errors: map[string]error{
				"users_repo_find":      nil,
				"jwt":                  nil,
				"uuid":                 errors.New("make jwt error"),
				"sessions_repo_delete": nil,
				"sessions_repo_save":   nil,
			},
			isPasswordValid: true,
			wantErr:         assert.Error,
		},
		{
			name: "delete existing session error",
			mockTimes: map[string]int{
				"users_repo_find":      1,
				"password_hash":        1,
				"jwt":                  1,
				"uuid":                 1,
				"get_refresh_ttl":      1,
				"sessions_repo_delete": 1,
				"sessions_repo_save":   0,
			},
			errors: map[string]error{
				"users_repo_find":      nil,
				"jwt":                  nil,
				"uuid":                 nil,
				"sessions_repo_delete": errors.New("delete session error"),
				"sessions_repo_save":   nil,
			},
			isPasswordValid: true,
			wantErr:         assert.Error,
		},
		{
			name: "save session error",
			mockTimes: map[string]int{
				"users_repo_find":      1,
				"password_hash":        1,
				"jwt":                  1,
				"uuid":                 1,
				"get_refresh_ttl":      1,
				"sessions_repo_delete": 1,
				"sessions_repo_save":   1,
			},
			errors: map[string]error{
				"users_repo_find":      nil,
				"jwt":                  nil,
				"uuid":                 nil,
				"sessions_repo_delete": nil,
				"sessions_repo_save":   errors.New("save session error"),
			},
			isPasswordValid: true,
			wantErr:         assert.Error,
		},
		{
			name: "success",
			mockTimes: map[string]int{
				"users_repo_find":      1,
				"password_hash":        1,
				"jwt":                  1,
				"uuid":                 1,
				"get_refresh_ttl":      1,
				"sessions_repo_delete": 1,
				"sessions_repo_save":   1,
			},
			errors: map[string]error{
				"users_repo_find":      nil,
				"jwt":                  nil,
				"uuid":                 nil,
				"sessions_repo_delete": nil,
				"sessions_repo_save":   nil,
			},
			isPasswordValid: true,
			wantErr:         assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			user := entities.User{
				ID:         userId,
				Email:      "email@mall.ru",
				Password:   "password_hash",
				VerifiedAt: sql.NullTime{Valid: true},
			}
			session := entities.Session{
				ID:          sessionId,
				UserID:      userId,
				FingerPrint: "unique_fingerprint",
				ExpiresAt:   refreshExpiration,
			}

			usersRepo.EXPECT().FindByEmail(gomock.Any(), "email@mall.ru").Return(&user, tt.errors["users_repo_find"]).Times(tt.mockTimes["users_repo_find"])
			hash.EXPECT().IsValidPasswordHash("password", "password_hash").Return(tt.isPasswordValid).Times(tt.mockTimes["password_hash"])
			jwtManager.EXPECT().MakeJWT(user.ID.String()).Return("jwt.token", tt.errors["jwt"]).Times(tt.mockTimes["jwt"])
			uuid.EXPECT().Generate().Return(sessionId, tt.errors["uuid"]).Times(tt.mockTimes["uuid"])
			jwtManager.EXPECT().MakeRefreshExpiration().Return(refreshExpiration).Times(tt.mockTimes["get_refresh_ttl"])
			sessionsRepo.EXPECT().DeleteByUserAndFingerPrint(gomock.Any(), user.ID.String(), "unique_fingerprint").Return(tt.errors["sessions_repo_delete"]).Times(tt.mockTimes["sessions_repo_delete"])
			sessionsRepo.EXPECT().Save(gomock.Any(), session).Return(tt.errors["sessions_repo_save"]).Times(tt.mockTimes["sessions_repo_save"])

			token, err := auth.Login(context.TODO(), "email@mall.ru", "password", "unique_fingerprint")

			tt.wantErr(t, err, fmt.Sprintf("Login()"))

			if token != nil {
				assert.Equal(t, token.RefreshToken, session.ID.String())
				assert.Equal(t, token.AccessToken, "jwt.token")
			}
		})
	}
}

func TestAuthUseCase_Logout(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessionsRepo := mocks.NewMockSessionsStorage(ctrl)

	auth := &AuthUseCase{sessionsRepo: sessionsRepo}

	var sessionId [16]byte
	copy(sessionId[:], "1843a7d7-1268-345b-bdb9-ga3a0e4b34e8")

	tests := []struct {
		name    string
		error   error
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "logout error",
			error:   errors.New("delete session error"),
			wantErr: assert.Error,
		},
		{
			name:    "logout success",
			error:   nil,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			session := entities.Session{ID: sessionId}

			sessionsRepo.EXPECT().DeleteByID(gomock.Any(), session.ID.String()).Return(tt.error).Times(1)
			tt.wantErr(t, auth.Logout(context.TODO(), session.ID.String()), fmt.Sprintf("Logout(ctx, %v)", session.ID.String()))
		})
	}
}

func TestAuthUseCase_Refresh(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uuid := mocks.NewMockUUIDGenerator(ctrl)
	jwtManager := mocks.NewMockJWTManager(ctrl)
	sessionsRepo := mocks.NewMockSessionsStorage(ctrl)

	auth := &AuthUseCase{
		uuid:         uuid,
		jwt:          jwtManager,
		sessionsRepo: sessionsRepo,
	}

	var userId [16]byte
	var sessionId [16]byte
	var newSessionId [16]byte
	copy(userId[:], "1955a7d6-0968-425b-bdb6-fb9a0e4b39e7")
	copy(sessionId[:], "1843a7d7-1268-345b-bdb9-ga3a0e4b34e8")
	copy(newSessionId[:], "7521z7d5-8791-245a-zda7-hj5a1e3b34c5")
	refreshExpiration := time.Now().Add(time.Second * 10)

	tests := []struct {
		name        string
		mockTimes   map[string]int
		errors      map[string]error
		fingerPrint string
		wantErr     assert.ErrorAssertionFunc
	}{
		{
			name: "find session error",
			mockTimes: map[string]int{
				"sessions_repo_find":   1,
				"sessions_repo_delete": 0,
				"jwt":                  0,
				"uuid":                 0,
				"get_refresh_ttl":      0,
				"sessions_repo_save":   0,
			},
			errors: map[string]error{
				"sessions_repo_find":   errors.New("find session error"),
				"sessions_repo_delete": nil,
				"jwt":                  nil,
				"uuid":                 nil,
				"sessions_repo_save":   nil,
			},
			fingerPrint: "unique_fingerprint",
			wantErr:     assert.Error,
		},
		{
			name: "delete session error",
			mockTimes: map[string]int{
				"sessions_repo_find":   1,
				"sessions_repo_delete": 1,
				"jwt":                  0,
				"uuid":                 0,
				"get_refresh_ttl":      0,
				"sessions_repo_save":   0,
			},
			errors: map[string]error{
				"sessions_repo_find":   nil,
				"sessions_repo_delete": errors.New("delete session error"),
				"jwt":                  nil,
				"uuid":                 nil,
				"sessions_repo_save":   nil,
			},
			fingerPrint: "unique_fingerprint",
			wantErr:     assert.Error,
		},
		{
			name: "invalid fingerPrint",
			mockTimes: map[string]int{
				"sessions_repo_find":   1,
				"sessions_repo_delete": 1,
				"jwt":                  0,
				"uuid":                 0,
				"get_refresh_ttl":      0,
				"sessions_repo_save":   0,
			},
			errors: map[string]error{
				"sessions_repo_find":   nil,
				"sessions_repo_delete": nil,
				"jwt":                  nil,
				"uuid":                 nil,
				"sessions_repo_save":   nil,
			},
			fingerPrint: "invalid_fingerprint",
			wantErr:     assert.Error,
		},
		{
			name: "make jwt error",
			mockTimes: map[string]int{
				"sessions_repo_find":   1,
				"sessions_repo_delete": 1,
				"jwt":                  1,
				"uuid":                 0,
				"get_refresh_ttl":      0,
				"sessions_repo_save":   0,
			},
			errors: map[string]error{
				"sessions_repo_find":   nil,
				"sessions_repo_delete": nil,
				"jwt":                  errors.New("make jwt error"),
				"uuid":                 nil,
				"sessions_repo_save":   nil,
			},
			fingerPrint: "unique_fingerprint",
			wantErr:     assert.Error,
		},
		{
			name: "new session uuid error",
			mockTimes: map[string]int{
				"sessions_repo_find":   1,
				"sessions_repo_delete": 1,
				"jwt":                  1,
				"uuid":                 1,
				"get_refresh_ttl":      0,
				"sessions_repo_save":   0,
			},
			errors: map[string]error{
				"sessions_repo_find":   nil,
				"sessions_repo_delete": nil,
				"jwt":                  nil,
				"uuid":                 errors.New("uuid error"),
				"sessions_repo_save":   nil,
			},
			fingerPrint: "unique_fingerprint",
			wantErr:     assert.Error,
		},
		{
			name: "save new session error",
			mockTimes: map[string]int{
				"sessions_repo_find":   1,
				"sessions_repo_delete": 1,
				"jwt":                  1,
				"uuid":                 1,
				"get_refresh_ttl":      1,
				"sessions_repo_save":   1,
			},
			errors: map[string]error{
				"sessions_repo_find":   nil,
				"sessions_repo_delete": nil,
				"jwt":                  nil,
				"uuid":                 nil,
				"sessions_repo_save":   errors.New("save session error"),
			},
			fingerPrint: "unique_fingerprint",
			wantErr:     assert.Error,
		},
		{
			name: "success",
			mockTimes: map[string]int{
				"sessions_repo_find":   1,
				"sessions_repo_delete": 1,
				"jwt":                  1,
				"uuid":                 1,
				"get_refresh_ttl":      1,
				"sessions_repo_save":   1,
			},
			errors: map[string]error{
				"sessions_repo_find":   nil,
				"sessions_repo_delete": nil,
				"jwt":                  nil,
				"uuid":                 nil,
				"sessions_repo_save":   nil,
			},
			fingerPrint: "unique_fingerprint",
			wantErr:     assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			session := entities.Session{
				ID:          sessionId,
				UserID:      userId,
				FingerPrint: "unique_fingerprint",
				ExpiresAt:   refreshExpiration,
			}

			newSession := entities.Session{
				ID:          newSessionId,
				UserID:      userId,
				FingerPrint: "unique_fingerprint",
				ExpiresAt:   refreshExpiration,
			}

			sessionsRepo.EXPECT().FindByID(gomock.Any(), session.ID.String()).Return(&session, tt.errors["sessions_repo_find"]).Times(tt.mockTimes["sessions_repo_find"])
			sessionsRepo.EXPECT().DeleteByID(gomock.Any(), session.ID.String()).Return(tt.errors["sessions_repo_delete"]).Times(tt.mockTimes["sessions_repo_delete"])
			jwtManager.EXPECT().MakeJWT(session.UserID.String()).Return("jwt.token", tt.errors["jwt"]).Times(tt.mockTimes["jwt"])
			uuid.EXPECT().Generate().Return(newSession.ID, tt.errors["uuid"]).Times(tt.mockTimes["uuid"])
			jwtManager.EXPECT().MakeRefreshExpiration().Return(refreshExpiration).Times(tt.mockTimes["get_refresh_ttl"])
			sessionsRepo.EXPECT().Save(gomock.Any(), newSession).Return(tt.errors["sessions_repo_save"]).Times(tt.mockTimes["sessions_repo_save"])

			token, err := auth.Refresh(context.TODO(), session.ID.String(), tt.fingerPrint)

			tt.wantErr(t, err, fmt.Sprintf("Login()"))

			if token != nil {
				assert.Equal(t, token.RefreshToken, newSession.ID.String())
				assert.Equal(t, token.AccessToken, "jwt.token")
			}
		})
	}
}
