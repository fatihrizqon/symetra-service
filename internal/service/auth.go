package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/fatihrizqon/symetra-service/internal/delivery/http/request"
	"github.com/fatihrizqon/symetra-service/internal/entity"
	"github.com/fatihrizqon/symetra-service/internal/repository"
	"github.com/fatihrizqon/symetra-service/internal/util"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type AuthResult struct {
	AccessToken  string
	RefreshToken string
	User         entity.User
}

type IAuthService interface {
	Login(req request.LoginRequest) (AuthResult, error)
	RefreshToken(refreshToken string) (AuthResult, error)
	Logout(refreshToken string) error
}

type AuthService struct {
	IAuthRepository  repository.IAuthRepository
	ITokenRepository repository.ITokenRepository
	validate         *validator.Validate
}

func NewAuthService(
	authRepo repository.IAuthRepository,
	tokenRepo repository.ITokenRepository,
	validate *validator.Validate,
) IAuthService {
	return &AuthService{
		IAuthRepository:  authRepo,
		ITokenRepository: tokenRepo,
		validate:         validate,
	}
}

func (e *AuthService) Login(req request.LoginRequest) (AuthResult, error) {
	var res AuthResult

	user, err := e.IAuthRepository.Login(req.Email)
	if err != nil {
		return res, err
	}

	if err := ValidatePassword(req.Password, user.Password); err != nil {
		return res, errors.New("credentials does not matches our record")
	}

	session := entity.Session{
		ID:     util.GenerateUUID(),
		UserID: user.Id,
	}

	session, err = e.ITokenRepository.CreateSession(session)
	if err != nil {
		return res, err
	}

	accessToken, err := util.CreateAccessToken(user, session.ID)
	if err != nil {
		return res, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := util.CreateRefreshToken(user, session.ID)
	if err != nil {
		return res, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	credential := entity.Credential{
		ID:           util.GenerateUUID(),
		SessionID:    session.ID,
		Type:         "REFRESH_TOKEN",
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(15 * time.Minute),
	}

	if err := e.ITokenRepository.CreateCredential(credential); err != nil {
		return res, err
	}

	return AuthResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

func (e *AuthService) RefreshToken(refreshToken string) (AuthResult, error) {
	oldCredential, err := e.ITokenRepository.FindCredentialByToken(refreshToken)
	if err != nil {
		return AuthResult{}, errors.New("invalid or revoked refresh token")
	}

	session, err := e.ITokenRepository.FindSessionByID(oldCredential.SessionID)
	if err != nil {
		return AuthResult{}, errors.New("session revoked or not found")
	}

	accessToken, err := util.CreateAccessToken(session.User, session.ID)
	if err != nil {
		return AuthResult{}, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, err := util.CreateRefreshToken(session.User, session.ID)
	if err != nil {
		return AuthResult{}, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	newCredential := entity.Credential{
		ID:           util.GenerateUUID(),
		SessionID:    session.ID,
		Type:         "REFRESH_TOKEN",
		RefreshToken: newRefreshToken,
		ExpiresAt:    time.Now().Add(15 * time.Minute),
	}

	if err := e.ITokenRepository.RevokeCredentialByID(oldCredential.ID); err != nil {
		return AuthResult{}, fmt.Errorf("failed to revoke old refresh token: %w", err)
	}

	if err := e.ITokenRepository.CreateCredential(newCredential); err != nil {
		return AuthResult{}, fmt.Errorf("failed to save new refresh token: %w", err)
	}

	return AuthResult{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		User:         session.User,
	}, nil
}

func (e *AuthService) Logout(refreshToken string) error {
	credential, err := e.ITokenRepository.FindCredentialByToken(refreshToken)
	if err != nil {
		return err
	}
	if err := e.ITokenRepository.RevokeCredentialByID(credential.ID); err != nil {
		return err
	}

	if err := e.ITokenRepository.RevokeSession(credential.SessionID); err != nil {
		return err
	}

	return nil
}

func ValidatePassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	)
}
