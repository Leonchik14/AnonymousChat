package service

import (
	"authentication-service/internal/repository"
	"context"
	"errors"
	"strconv"
	"time"

	"authentication-service/pkg/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo     *repository.UserRepository
	redisRepo    *repository.RedisRepository
	emailService *EmailService
	jwtSecret    string
	accessTTL    time.Duration
	refreshTTL   time.Duration
}

func NewAuthService(userRepo *repository.UserRepository, redisRepo *repository.RedisRepository, emailService *EmailService, jwtSecret string, accessTTL time.Duration, refreshTTL time.Duration) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		redisRepo:    redisRepo,
		emailService: emailService,
		jwtSecret:    jwtSecret,
		accessTTL:    accessTTL,
		refreshTTL:   refreshTTL,
	}
}

// RegisterUser - регистрация пользователя + отправка email
func (s *AuthService) RegisterUser(ctx context.Context, email, password string) error {
	existingUser, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("пользователь уже существует")
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		Email:        email,
		PasswordHash: string(passwordHash),
	}
	return s.userRepo.CreateUser(ctx, user)
}

func (s *AuthService) RequestEmailVerification(ctx context.Context, email string) error {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("пользователь не найден")
	}
	if user.IsVerified {
		return errors.New("email уже подтверждён")
	}

	token, err := s.generateJWT(int64(user.ID), s.accessTTL)
	if err != nil {
		return err
	}
	return s.emailService.SendVerificationEmail(email, token)
}

func (s *AuthService) LoginUser(ctx context.Context, email, password string) (string, string, uint, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil || user == nil {
		return "", "", 0, errors.New("неверный email или пароль")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", 0, errors.New("неверный email или пароль")
	}

	accessToken, err := s.generateJWT(int64(user.ID), s.accessTTL)
	if err != nil {
		return "", "", 0, err
	}
	refreshToken, err := s.generateJWT(int64(user.ID), s.refreshTTL)
	if err != nil {
		return "", "", 0, err
	}

	err = s.redisRepo.SetSession(ctx, refreshToken, user.ID, s.refreshTTL)
	if err != nil {
		return "", "", 0, err
	}

	return accessToken, refreshToken, user.ID, nil
}

func (s *AuthService) LogoutUser(ctx context.Context, token string) error {
	return s.redisRepo.DeleteSession(ctx, token)
}

func (s *AuthService) ValidateToken(ctx context.Context, token string) (int64, error) {
	claims := &jwt.RegisteredClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})
	if err != nil || !parsedToken.Valid {
		return 0, errors.New("недействительный токен")
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return 0, errors.New("токен истек")
	}

	userID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return 0, errors.New("неверный формат userID в токене")
	}

	return userID, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	userID, err := s.redisRepo.GetUserIDBySession(ctx, refreshToken)
	if err != nil || userID == 0 {
		return "", errors.New("refresh token недействителен, выполните повторный вход")
	}

	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return "", errors.New("пользователь не найден")
	}

	if !user.IsVerified {
		ttl, err := s.redisRepo.GetSessionTTL(ctx, refreshToken)
		if err != nil || ttl <= 0 {
			return "", errors.New("сессия истекла, выполните повторный вход")
		}
	}

	accessToken, err := s.generateJWT(userID, s.accessTTL)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (s *AuthService) VerifyUser(ctx context.Context, userID int64) error {
	return s.userRepo.SetUserVerified(ctx, userID)
}

func (s *AuthService) generateJWT(userID int64, ttl time.Duration) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(userID, 10),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
