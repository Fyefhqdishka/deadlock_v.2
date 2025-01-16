package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Fyefhqdishka/deadlock_v.2/internal/models"
	"github.com/google/uuid"
	"log/slog"
)

type Storage struct {
	db  PgxPoolIface
	log *slog.Logger
}

func (s *Storage) Create(ctx context.Context, user *models.User) error {
	s.log.Debug("starting registration")

	stmt := `
		INSERT INTO users (name, username, email, password, gender, dob, avatar) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	Avatar := "https://i.ibb.co.com/LNYxmRs/pngtree-character-default-avatar-image-2237203.jpg"

	user.Avatar = Avatar

	res, err := s.db.Exec(ctx, stmt,
		user.Name,
		user.Username,
		user.Email,
		user.Password,
		user.Gender,
		user.Dob,
		user.Avatar,
	)

	s.log.Debug("SQL query:", stmt)

	s.log.Debug("registration ", res)

	if err != nil {
		s.log.Debug("err:", err, "data:", res)
		return err
	}

	s.log.Debug("registration successfully", "username:", user.Username)

	return nil
}

func (s *Storage) Login(ctx context.Context, username, password string) (string, error) {
	s.log.Debug("LoginUser", "starting login user")

	var UserID string
	var passwordHash string

	s.log.Debug("Login", "username", username)
	stmt := `SELECT id, password FROM users WHERE username = $1`
	err := s.db.QueryRow(ctx, stmt, username).Scan(&UserID, &passwordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			s.log.Warn(
				"LoginUser",
				"пользователь не найден",
				"username", username,
			)
			return "", fmt.Errorf("пользователь не найден")
		}
		s.log.Error(
			"LoginUser",
			"ошибка выполнения SQL-запроса при логине",
			"err", err,
		)
		return "", err
	}

	if !models.CheckPasswordHash(password, passwordHash) {
		s.log.Warn("LoginUser", "неверный пароль для пользователя", "username", username)
		return "", fmt.Errorf("неверный пароль")
	}

	sessionID := uuid.New().String()

	stmt = `INSERT INTO sessions (user_id, session_id) VALUES ($1, $2)`
	_, err = s.db.Exec(ctx, stmt, UserID, sessionID)
	if err != nil {
		s.log.Error("LoginUser", "ошибка сохранения сессии", "err", err)
		return "", err
	}

	s.log.Info("пользователь успешно аутентифицирован", "username", username)

	return sessionID, nil
}

func (s *Storage) GetAllUsers() ([]models.User, error) {
	var users []models.User

	rows, err := s.db.Query(context.Background(), "SELECT id, name, username, email, password, gender, dob, avatar, time_registration FROM users")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.UserID, &user.Name, &user.Username, &user.Email, &user.Password, &user.Gender, &user.Dob, &user.Avatar, &user.Time_registration); err != nil {
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating over rows: %v", err)
	}

	return users, nil
}

func (s *Storage) CreateDialog(ctx context.Context, userIDOne, userIDTwo string) (int64, error) {
	var dialogID int64

	stmt := `INSERT INTO dialogs (user_id_1, user_id_2) 
	         VALUES ($1, $2) 
	         ON CONFLICT (user_id_1, user_id_2) DO NOTHING 
	         RETURNING id`

	err := s.db.QueryRow(ctx, stmt, userIDOne, userIDTwo).Scan(&dialogID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return 0, nil
		}

		s.log.Error(
			"CreateDialog",
			"Ошибка добавления диалога в БД",
			"err:", err,
		)

		return 0, err
	}

	return dialogID, nil
}

func (s *Storage) GetUserByID(ctx context.Context, sessionID string) (string, error) {
	if sessionID == "" {
		fmt.Println("Session ID is missing")
		return "", fmt.Errorf("Session ID is missing")
	}

	var userID string
	stmt := `SELECT user_id FROM sessions WHERE session_id = $1`
	err := s.db.QueryRow(ctx, stmt, sessionID).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			s.log.Error("Session ID not found")
			return "", fmt.Errorf("Session not found for ID: %s", sessionID)
		}
		return "", fmt.Errorf("Error querying session: %v", err)
	}

	return userID, nil
}

func (s *Storage) GetUserDialogs(ctx context.Context, userID string) ([]models.Dialog, error) {
	var dialogs []models.Dialog

	stmt := `SELECT 
    d.id AS dialog_id, 
    CASE 
        WHEN d.user_id_1 = $1 THEN u2.username 
        ELSE u1.username 
    END AS opponent_username,
    CASE 
        WHEN d.user_id_1 = $1 THEN u2.avatar 
        ELSE u1.avatar 
    END AS opponent_avatar
FROM dialogs d
JOIN users u1 ON u1.id = d.user_id_1 
JOIN users u2 ON u2.id = d.user_id_2 
WHERE d.user_id_1 = $1 OR d.user_id_2 = $1;
`

	rows, err := s.db.Query(ctx, stmt, userID)
	if err != nil {
		s.log.Error("GetUserDialogs", "Failed to fetch user dialogs", "err:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var dialog models.Dialog
		if err := rows.Scan(&dialog.DialogID, &dialog.UserTwoUsername, &dialog.Avatar); err != nil {
			s.log.Error("GetUserDialogs", "Failed to scan dialog row", "err:", err)
			return nil, err
		}
		dialogs = append(dialogs, dialog)
	}

	if err := rows.Err(); err != nil {
		s.log.Error("GetUserDialogs", "Error reading rows", "err:", err)
		return nil, err
	}

	return dialogs, nil
}

func (s *Storage) CreatePost(ctx context.Context, post *models.Post, userID string) error {
	s.log.Debug("Обработка запроса на добавление поста в БД")

	stmt := `INSERT INTO Posts (title, content, user_id) VALUES ($1, $2, $3)`
	_, err := s.db.Exec(ctx, stmt, post.Title, post.Content, userID)
	if err != nil {
		s.log.Error(
			"PostCreate",
			"Ошибка при добавлении нового поста",
		)
		return err
	}

	s.log.Debug(
		"PostCreate",
		"Добавление поста в БД прошло успешно",
	)

	return nil
}

func (s *Storage) GetAllPosts(ctx context.Context) ([]models.Post, error) {
	s.log.Debug("Получение всех постов")

	stmt := `SELECT 
    posts.id,
    posts.title,
    posts.content,
    TO_CHAR(posts.created_at, 'YYYY-MM-DD HH24:MI:SS') AS created_at,
    users.username,
    users.avatar
FROM posts
JOIN users ON posts.user_id = users.id;
`

	rows, err := s.db.Query(ctx, stmt)
	if err != nil {
		s.log.Error("GetAllPosts", "Failed to fetch posts", "err:", err)
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post

	for rows.Next() {
		var post models.Post
		var createdAt string
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &createdAt, &post.Username, &post.Avatar); err != nil {
			s.log.Error("GetAllPosts", "Failed to scan post row", "err:", err)
			return nil, err
		}
		post.CreateAt = createdAt
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		s.log.Error("GetAllPosts", "Error reading rows", "err:", err)
		return nil, err
	}

	s.log.Debug(
		"Все посты получены успешно",
	)

	return posts, nil
}

func (s *Storage) GetPost(ctx context.Context, post *models.Post) error {
	s.log.Info("Fetching post", "postID", post.ID)

	stmt := `SELECT p.title, p.content, TO_CHAR(p.created_at, 'YYYY-MM-DD HH24:MI:SS') AS created_at, u.username, u.avatar
	FROM posts p
	JOIN users u ON p.user_id = u.id
	WHERE p.id = $1;`

	s.log.Info("Executing query", "query", stmt, "postID", post.ID)

	err := s.db.QueryRow(ctx, stmt, post.ID).Scan(&post.Title, &post.Content, &post.CreateAt, &post.Username, &post.Avatar)
	if err != nil {
		if err == sql.ErrNoRows {
			s.log.Error("Post not found", "postID", post.ID)
			return fmt.Errorf("post with ID %d not found", post.ID)
		}
		s.log.Error("Failed to get post", "postID", post.ID, "err", err)
		return fmt.Errorf("failed to get post with ID %d: %v", post.ID, err)
	}

	s.log.Info("Successfully fetched post", "postID", post.ID, "title", post.Title)

	return nil
}
