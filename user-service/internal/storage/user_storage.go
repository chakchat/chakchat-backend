package storage

import (
	"context"
	"errors"
	"time"

	"github.com/chakchat/chakchat-backend/user-service/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("not found")
var ErrAlreadyExists = errors.New("already exists")

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
	db *gorm.DB
}

func NewUserStorage(db *gorm.DB) *UserStorage {
	return &UserStorage{db: db}
}

func (s *UserStorage) GetUserByPhone(ctx context.Context, phone string) (*models.User, error) {
	var user models.User
	if err := s.db.WithContext(ctx).Where(&models.User{Phone: phone}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (s *UserStorage) GetUserById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := s.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (s *UserStorage) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	if err := s.db.WithContext(ctx).Where(&models.User{Username: username}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (s *UserStorage) GetUsersByCriteria(ctx context.Context, req SearchUsersRequest) (*SearchUsersResponse, error) {
	var users []models.User
	query := s.db.WithContext(ctx).Model(&users)

	if req.Name != nil {
		query = query.Where(&models.User{Name: *req.Name})
		if err := query.Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrNotFound
			}
			return nil, err
		}
	}

	if req.Username != nil {
		query = query.Where(&models.User{Username: *req.Username})
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
		Users:  users,
		Offset: offset/limit + 1,
		Count:  int(count),
	}, nil
}

func (s *UserStorage) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	var newUser models.User = models.User{
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

func (s *UserStorage) UpdateUser(ctx context.Context, user *models.User, name string, username string, birthday *time.Time) (*models.User, error) {

	if err := s.db.WithContext(ctx).Save(&models.User{ID: (*user).ID, Name: name, Username: username, DateOfBirth: birthday}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return user, nil
}
