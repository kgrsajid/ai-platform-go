package userservice

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"project-go/internal/lib/password"
	"project-go/internal/lib/random"
	"project-go/internal/models"
)

type UserRepository interface {
	CreateUser(u *models.User) (*models.User, error)
	FindUserByEmail(email string) (*models.User, error)
	UpdatePassword(id uint, hash string) error
}

type PasswordResetRepository interface {
	Create(pr *models.PasswordReset) error
	FindByEmailAndCode(email, code string) (*models.PasswordReset, error)
	FindByToken(token string) (*models.PasswordReset, error)
	UpdateToken(id uint, token string) error
	DeleteByEmail(email string) error
}

type EmailSender interface {
	SendResetCode(to, code string) error
}

type Service struct {
	userRepo    UserRepository
	resetRepo   PasswordResetRepository
	emailSender EmailSender
}

func New(userRepo UserRepository, resetRepo PasswordResetRepository, emailSender EmailSender) *Service {
	return &Service{
		userRepo:    userRepo,
		resetRepo:   resetRepo,
		emailSender: emailSender,
	}
}

func (s *Service) CreateUser(user *models.User) (*models.User, error) {
	return s.userRepo.CreateUser(user)
}

func (s *Service) FindUserByEmail(email string) (*models.User, error) {
	return s.userRepo.FindUserByEmail(email)
}

func (s *Service) ForgotPassword(email string) error {
	if _, err := s.userRepo.FindUserByEmail(email); err != nil {
		return errors.New("user not found")
	}

	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return err
	}
	code := fmt.Sprintf("%06d", n.Int64())

	_ = s.resetRepo.DeleteByEmail(email)

	pr := &models.PasswordReset{
		Email:     email,
		Code:      code,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}
	if err := s.resetRepo.Create(pr); err != nil {
		return err
	}

	return s.emailSender.SendResetCode(email, code)
}

func (s *Service) VerifyCode(email, code string) (string, error) {
	pr, err := s.resetRepo.FindByEmailAndCode(email, code)
	if err != nil {
		return "", errors.New("invalid code")
	}

	if time.Now().After(pr.ExpiresAt) {
		return "", errors.New("code expired")
	}

	token := random.NewRandomString(32)
	if err := s.resetRepo.UpdateToken(pr.ID, token); err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) ResetPassword(token, newPassword string) error {
	pr, err := s.resetRepo.FindByToken(token)
	if err != nil {
		return errors.New("invalid token")
	}

	if time.Now().After(pr.ExpiresAt) {
		return errors.New("token expired")
	}

	user, err := s.userRepo.FindUserByEmail(pr.Email)
	if err != nil {
		return err
	}

	hash, err := password.HashPassword(newPassword)
	if err != nil {
		return err
	}

	if err := s.userRepo.UpdatePassword(user.ID, hash); err != nil {
		return err
	}

	return s.resetRepo.DeleteByEmail(pr.Email)
}
