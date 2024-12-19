package storage

import (
	"context"
	"github.com/Fyefhqdishka/deadlock_v.2/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage interface {
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
