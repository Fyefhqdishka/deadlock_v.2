package service

import (
	"context"
	"fmt"
	"github.com/Fyefhqdishka/deadlock_v.2/internal/models"
	"github.com/Fyefhqdishka/deadlock_v.2/internal/storage"
	"log/slog"
)

type ServiceIface interface {
	Create(ctx context.Context, user *models.User) error
	Login(ctx context.Context, username, password string) (string, error)
	GetAllUsers() ([]models.User, error)
	CreateDialog(ctx context.Context, userIDOne, userIDTwo string) (int64, error)
	GetUserByID(ctx context.Context, sessionID string) (string, error)
	GetUserDialogs(ctx context.Context, userID string) ([]models.Dialog, error)
	CreatePost(ctx context.Context, post *models.Post, userID string) error
	GetAllPosts(ctx context.Context) ([]models.Post, error)
	GetPost(ctx context.Context, post *models.Post) error
}

type Service struct {
	repo storage.Storage
	log  *slog.Logger
}

func NewService(repo storage.Storage, log *slog.Logger) *Service {
	return &Service{
		repo: repo,
		log:  log,
	}
}

func (s *Service) Create(ctx context.Context, user *models.User) error {
	hashedPassword, err := models.HashPassword(user.Password)
	if err != nil {
		fmt.Printf("ошибка при хэшировании пароля: %v\n", err)
		return err
	}
	user.Password = hashedPassword

	err = s.repo.Create(ctx, user)
	if err != nil {
		return err
	}

	return nil
}
func (s *Service) Login(ctx context.Context, username, password string) (string, error) {
	s.log.Debug("Login", "starting login attempt", "username", username)

	sessionID, err := s.repo.Login(ctx, username, password)
	if err != nil {
		s.log.Error("Login", "error during login", "username", username, "error", err)
		return "", fmt.Errorf("не удалось выполнить вход: %v", err)
	}

	s.log.Info("Login", "user authenticated successfully", "username", username)

	return sessionID, nil
}

func (s *Service) GetAllUsers() ([]models.User, error) {
	res, err := s.repo.GetAllUsers()
	if err != nil {
		return []models.User{}, err
	}

	return res, nil
}

func (s *Service) CreateDialog(ctx context.Context, userIDOne, userIDTwo string) (int64, error) {
	dialogID, err := s.repo.CreateDialog(ctx, userIDOne, userIDTwo)
	if err != nil {
		return 0, err
	}

	return dialogID, nil
}

func (s *Service) GetUserByID(ctx context.Context, sessionID string) (string, error) {
	userID, err := s.repo.GetUserByID(ctx, sessionID)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func (s *Service) GetUserDialogs(ctx context.Context, userID string) ([]models.Dialog, error) {
	users, err := s.repo.GetUserDialogs(ctx, userID)
	if err != nil {
		return []models.Dialog{}, err
	}

	return users, nil
}

func (s *Service) CreatePost(ctx context.Context, post *models.Post, userID string) error {
	err := s.repo.CreatePost(ctx, post, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetAllPosts(ctx context.Context) ([]models.Post, error) {
	posts, err := s.repo.GetAllPosts(ctx)
	if err != nil {
		return []models.Post{}, err
	}

	return posts, nil
}

func (s *Service) GetPost(ctx context.Context, post *models.Post) error {
	err := s.repo.GetPost(ctx, post)
	if err != nil {
		return err
	}

	return nil
}
