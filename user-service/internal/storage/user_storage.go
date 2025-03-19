package storage

import (
	"context"
	"errors"
	"time"

	"github.com/chakchat/chakchat-backend/user-service/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("not found")
var ErrAlreadyExists = errors.New("already exists")

type UpdateUserRequest struct {
	Name        string
	Username    string
	DateOfBirth *time.Time
}

type SearchUsersRequest struct {
	Name     *string
	Username *string
	Offset   *int
	Limit    *int
}

type SearchUsersResponse struct {
	Users  []models.User
	Offset int
	Count  int
}

type UserStorage struct {
	db *pgxpool.Pool
}

func NewUserStorage(db *pgxpool.Pool) *UserStorage {
	return &UserStorage{db: db}
}

func (s *UserStorage) GetUserByPhone(ctx context.Context, phone string) (*models.User, error) {
	var user models.User
	q := `SELECT
		id,
		name,
		username,
		phone,
		date_of_birth,
		photo_url,
		created_at,
		date_of_birth_visibility,
		phone_visibility 
	FROM users
	WHERE phone = $1`

	row := s.db.QueryRow(ctx, q, phone)
	if err := row.Scan(&user.ID, &user.Name, &user.Username, &user.Phone,
		&user.DateOfBirth, &user.PhotoURL, &user.CreatedAt,
		&user.DateOfBirthVisibility, &user.PhoneVisibility); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (s *UserStorage) GetUserById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	q := `SELECT
		id,
		name,
		username,
		phone,
		date_of_birth,
		photo_url,
		created_at,
		date_of_birth_visibility,
		phone_visibility 
	FROM users
	WHERE id = $1`

	row := s.db.QueryRow(ctx, q, id)
	if err := row.Scan(&user.ID, &user.Name, &user.Username, &user.Phone,
		&user.DateOfBirth, &user.PhotoURL, &user.CreatedAt,
		&user.DateOfBirthVisibility, &user.PhoneVisibility); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (s *UserStorage) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	q := `SELECT
		id,
		name,
		username,
		phone,
		date_of_birth,
		photo_url,
		created_at,
		date_of_birth_visibility,
		phone_visibility 
	FROM users
	WHERE username = $1`

	row := s.db.QueryRow(ctx, q, username)
	if err := row.Scan(&user.ID, &user.Name, &user.Username, &user.Phone,
		&user.DateOfBirth, &user.PhotoURL, &user.CreatedAt,
		&user.DateOfBirthVisibility, &user.PhoneVisibility); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (s *UserStorage) GetUsersByCriteria(ctx context.Context, req SearchUsersRequest) (*SearchUsersResponse, error) {
	var users []models.User
	q := `SELECT
		id,
		name,
		username,
		phone,
		date_of_birth,
		photo_url,
		created_at,
		date_of_birth_visibility,
		phone_visibility 
	FROM users
	WHERE 1=1`

	args := []interface{}{}
	if req.Name != nil {
		q += ` AND name ILIKE $1`
		args = append(args, "%"+*req.Name+"%")
	}

	if req.Username != nil {
		q += ` AND username ILIKE $2`
		args = append(args, "%"+*req.Username+"%")
	}

	offset := 0
	if req.Offset != nil {
		offset = *req.Offset
	}

	limit := 10
	if req.Limit != nil {
		limit = *req.Limit
	}

	countQuery := `SELECT COUNT(*) FROM users WHERE 1=1`
	if req.Name != nil {
		countQuery += ` AND name ILIKE $1`
	}
	if req.Username != nil {
		countQuery += ` AND username ILIKE $2`
	}

	var count int64
	if err := s.db.QueryRow(ctx, countQuery, args...).Scan(&count); err != nil {
		return nil, err
	}
	count -= int64(offset)

	q += ` LIMIT $3 OFFSET $4`
	args = append(args, limit, offset)

	rows, err := s.db.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.Name, &user.Phone, &user.DateOfBirth, &user.PhotoURL, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return &SearchUsersResponse{
		Users:  users,
		Offset: offset,
		Count:  int(count),
	}, nil
}

func (s *UserStorage) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {

	q := `SELECT id FROM users WHERE username = $1 OR phone = $2`
	var existingID uuid.UUID
	err := s.db.QueryRow(ctx, q, user.Username, user.Phone).Scan(&existingID)
	if err == nil {
		return nil, ErrAlreadyExists
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	var newUser models.User = models.User{
		ID:          uuid.New(),
		Username:    user.Username,
		Name:        user.Name,
		Phone:       user.Phone,
		DateOfBirth: user.DateOfBirth,
		PhotoURL:    user.PhotoURL,
		CreatedAt:   time.Now().Unix(),
	}

	insertQuery := `INSERT INTO users (id, username, name, phone, date_of_birth, photo_url, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = s.db.Exec(ctx, insertQuery, newUser.ID, newUser.Username, newUser.Name, newUser.Phone, newUser.DateOfBirth, newUser.PhotoURL, newUser.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

func (s *UserStorage) UpdateUser(ctx context.Context, user *models.User, req *UpdateUserRequest) (*models.User, error) {

	q := `SELECT id FROM users WHERE username = $1`
	var existingId uuid.UUID
	err := s.db.QueryRow(ctx, q, req.Username).Scan(&existingId)
	if err == nil && existingId != user.ID {
		return nil, ErrAlreadyExists
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	updateQuery := `UPDATE users SET name = $1, username = $2, date_of_birth = $3 WHERE id = $4`
	_, err = s.db.Exec(ctx, updateQuery, req.Name, req.Username, req.DateOfBirth, user.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	updatedUser, err := s.GetUserById(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	return updatedUser, nil
}

func (s *UserStorage) UpdatePhoto(ctx context.Context, id uuid.UUID, photoURL string) (*models.User, error) {
	updateQuery := `UPDATE users SET photo_url = $1 WHERE id = $2`
	result, err := s.db.Exec(ctx, updateQuery, photoURL, id)
	if err != nil {
		return nil, err
	}
	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return nil, ErrNotFound
	}

	user, err := s.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserStorage) DeletePhoto(ctx context.Context, id uuid.UUID) (*models.User, error) {
	updateQuery := `UPDATE users SET photo_url = '' WHERE id = $1`
	_, err := s.db.Exec(ctx, updateQuery, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	user, err := s.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
