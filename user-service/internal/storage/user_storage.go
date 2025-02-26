package storage

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("not found")
var ErrAlreadyExists = errors.New("already exists")

type User struct {
	ID          uuid.UUID  `gorm:"primaryKey" json:"id"`
	Username    string     `json:"name"`
	Name        string     `json:"username"`
	Phone       *string    `json:"phone,omitempty"`
	DateOfBirth *time.Time `json:"dateOfBirth,omitempty"`
	PhotoURL    string     `json:"photo"`
	CreatedAt   int64
}

type SearchUsersRequest struct {
	Name     *string
	Username *string
	Offset   *int
	Limit    *int
}

type SearchUsersResponse struct {
	Users []User `json:"users"`
	Page  int
	Count int
}

type UserStorage struct {
	db *gorm.DB
}

func NewUserStorage(db *gorm.DB) *UserStorage {
	return &UserStorage{db: db}
}

func (s *UserStorage) GetUserByPhone(ctx context.Context, phone string) (*User, error) {
	var user User
	if err := s.db.WithContext(ctx).Where(&User{Phone: &phone}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &User{
		ID:          user.ID,
		Username:    user.Username,
		Name:        user.Name,
		Phone:       user.Phone,
		DateOfBirth: user.DateOfBirth,
		PhotoURL:    user.PhotoURL,
		CreatedAt:   user.CreatedAt,
	}, nil
}

func (s *UserStorage) GetUserById(ctx context.Context, id uuid.UUID) (*User, error) {
	var user User
	if err := s.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &User{
		ID:          user.ID,
		Username:    user.Username,
		Name:        user.Name,
		Phone:       user.Phone,
		DateOfBirth: user.DateOfBirth,
		PhotoURL:    user.PhotoURL,
		CreatedAt:   user.CreatedAt,
	}, nil
}

func (s *UserStorage) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	if err := s.db.WithContext(ctx).Where(&User{Username: username}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &User{
		ID:          user.ID,
		Username:    user.Username,
		Name:        user.Name,
		Phone:       user.Phone,
		DateOfBirth: user.DateOfBirth,
		PhotoURL:    user.PhotoURL,
		CreatedAt:   user.CreatedAt,
	}, nil
}

func (s *UserStorage) GetUsersByCriteria(ctx context.Context, req SearchUsersRequest) (*SearchUsersResponse, error) {
	var users []User
	query := s.db.WithContext(ctx).Model(&users)

	if req.Name != nil {
		query = query.Where(&User{Name: *req.Name})
		if err := query.Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrNotFound
			}
			return nil, err
		}
	}

	if req.Username != nil {
		query = query.Where(&User{Username: *req.Username})
		if err := query.Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrNotFound
			}
			return nil, err
		}
	}

	offset := 0
	if req.Offset != nil {
		offset = *req.Offset
	}

	limit := 10
	if req.Limit != nil {
		limit = *req.Limit
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, err
	}

	if err := query.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, err
	}

	return &SearchUsersResponse{
		Users: users,
		Page:  offset/limit + 1,
		Count: int(count),
	}, nil
}

func (s *UserStorage) CreateUser(ctx context.Context, user *User) (*User, error) {
	var newUser User = User{
		ID:          uuid.New(),
		Username:    user.Username,
		Name:        user.Name,
		Phone:       user.Phone,
		DateOfBirth: user.DateOfBirth,
		PhotoURL:    user.PhotoURL,
		CreatedAt:   time.Now().Unix(),
	}

	if err := s.db.WithContext(ctx).First(&user).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	} else {
		return nil, ErrAlreadyExists
	}

	if err := s.db.Create(&newUser).Error; err != nil {
		return nil, err
	}
	return &newUser, nil
}
