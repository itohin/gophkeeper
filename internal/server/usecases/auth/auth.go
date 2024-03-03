package auth

import (
	"context"
	"errors"
	"github.com/itohin/gophkeeper/internal/server/entities"
	"time"
)

type Mailer interface {
	SendMail(to []string, message string) error
	SendMailAsync(to []string, message string)
}

type UsersStorage interface {
	Save(ctx context.Context, user entities.User) error
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	FindByID(ctx context.Context, id string) (*entities.User, error)
}

type SessionsStorage interface {
	Save(ctx context.Context, user entities.Session) error
	FindByID(ctx context.Context, id string) (*entities.Session, error)
	FindByFingerPrint(ctx context.Context, userId, fingerPrint string) (*entities.Session, error)
	DeleteByID(ctx context.Context, sessionID string) error
}

type PasswordHasher interface {
	HashPassword(password string) (string, error)
	IsValidPasswordHash(password, hash string) bool
}

type OTPGenerator interface {
	RandomSecret() (string, error)
}

type UUIDGenerator interface {
	Generate() ([16]byte, error)
}

type JWTManager interface {
	MakeJWT(userID string) (string, error)
	MakeRefreshExpiration() time.Time
}

type DBTransactionManager interface {
	Transaction(context.Context, func() error) error
}

type AuthUseCase struct {
	hash         PasswordHasher
	uuid         UUIDGenerator
	otp          OTPGenerator
	usersRepo    UsersStorage
	sessionsRepo SessionsStorage
	mailer       Mailer
	jwt          JWTManager
	tx           DBTransactionManager
}

func NewAuthUseCase(
	hash PasswordHasher,
	uuid UUIDGenerator,
	otp OTPGenerator,
	usersRepo UsersStorage,
	sessionsRepo SessionsStorage,
	mailer Mailer,
	jwt JWTManager,
	tx DBTransactionManager,
) *AuthUseCase {
	return &AuthUseCase{
		hash:         hash,
		uuid:         uuid,
		otp:          otp,
		usersRepo:    usersRepo,
		sessionsRepo: sessionsRepo,
		mailer:       mailer,
		jwt:          jwt,
		tx:           tx,
	}
}

func (a *AuthUseCase) Login(ctx context.Context, email, password, fingerprint string) (*entities.Token, error) {
	user, err := a.usersRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if !user.IsVerified() {
		return nil, errors.New("user not verified")
	}

	if !a.hash.IsValidPasswordHash(user.Password, password) {
		return nil, errors.New("wrong credentials")
	}
	accessToken, err := a.jwt.MakeJWT(user.ID.String())
	if err != nil {
		return nil, err
	}
	//TODO: проверки наличие активной сессии для данного устройства
	sessionID, err := a.uuid.Generate()
	if err != nil {
		return nil, err
	}
	//TODO: проверки безопасности: количество сессий у юзера(не более 5 устройств)...
	session := entities.NewSession(sessionID, user.ID, fingerprint, a.jwt.MakeRefreshExpiration())
	err = a.sessionsRepo.Save(ctx, *session)
	if err != nil {
		return nil, err
	}

	return entities.NewToken(accessToken, session.ID.String()), nil
}

func (a *AuthUseCase) Register(ctx context.Context, email, password string) error {

	passwordHash, err := a.hash.HashPassword(password)
	if err != nil {
		return err
	}
	uuid, err := a.uuid.Generate()
	if err != nil {
		return err
	}

	otp, err := a.otp.RandomSecret()
	if err != nil {
		return err
	}
	user := entities.NewUser(uuid, email, passwordHash, otp)
	err = a.usersRepo.Save(ctx, *user)
	if err != nil {
		return err
	}

	a.mailer.SendMailAsync(
		[]string{email},
		"Для подтверждения адреса электронной почты в сервисе gophkeeper, введите пожалуйста код подтверждения: "+otp,
	)

	return nil
}

func (a *AuthUseCase) Verify(ctx context.Context, email, otp, fingerprint string) (*entities.Token, error) {
	user, err := a.usersRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	err = user.Verify(otp)
	if err != nil {
		return nil, err
	}

	accessToken, err := a.jwt.MakeJWT(user.ID.String())
	if err != nil {
		return nil, err
	}
	sessionID, err := a.uuid.Generate()
	if err != nil {
		return nil, err
	}
	session := entities.NewSession(sessionID, user.ID, fingerprint, a.jwt.MakeRefreshExpiration())

	err = a.tx.Transaction(ctx, func() error {
		err := a.sessionsRepo.Save(ctx, *session)
		if err != nil {
			return err
		}
		err = a.usersRepo.Save(ctx, *user)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return entities.NewToken(accessToken, session.ID.String()), nil
}

func (a *AuthUseCase) Logout(ctx context.Context, sessionID string) error {
	err := a.sessionsRepo.DeleteByID(ctx, sessionID)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthUseCase) Refresh(ctx context.Context, sessionID, fingerprint string) (*entities.Token, error) {
	session, err := a.sessionsRepo.FindByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	err = a.sessionsRepo.DeleteByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if session.IsExpired() {
		return nil, errors.New("token expired")
	}
	if session.FingerPrint != fingerprint {
		return nil, errors.New("invalid token")
	}

	accessToken, err := a.jwt.MakeJWT(session.UserID.String())
	if err != nil {
		return nil, err
	}
	newSessionID, err := a.uuid.Generate()
	if err != nil {
		return nil, err
	}
	newSession := entities.NewSession(newSessionID, session.UserID, fingerprint, a.jwt.MakeRefreshExpiration())
	err = a.sessionsRepo.Save(ctx, *newSession)
	if err != nil {
		return nil, err
	}

	return entities.NewToken(accessToken, newSession.ID.String()), nil
}
