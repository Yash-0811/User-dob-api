package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/yash/user-dob-api/internal/models"
	"github.com/yash/user-dob-api/internal/repository"
)

const dobLayout = "2006-01-02"

// ErrNotFound is returned when a requested user does not exist.
var ErrNotFound = errors.New("user not found")

// UserService defines the business-logic contract.
type UserService interface {
	Create(ctx context.Context, req models.CreateUserRequest) (models.UserResponse, error)
	GetByID(ctx context.Context, id int32) (models.UserDetailResponse, error)
	Update(ctx context.Context, id int32, req models.UpdateUserRequest) (models.UserResponse, error)
	Delete(ctx context.Context, id int32) error
	List(ctx context.Context, page, pageSize int) (models.PaginatedUsersResponse, error)
}

type userService struct {
	repo repository.UserRepository
}

// New returns a UserService backed by the given repository.
func New(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

// CalculateAge returns the completed years between dob and now.
// Exported so it can be unit-tested directly.
func CalculateAge(dob time.Time) int {
	now := time.Now()
	years := now.Year() - dob.Year()
	// Step back one year if the birthday hasn't occurred yet this year.
	anniversary := dob.AddDate(years, 0, 0)
	if now.Before(anniversary) {
		years--
	}
	return years
}

func (s *userService) Create(ctx context.Context, req models.CreateUserRequest) (models.UserResponse, error) {
	dob, err := time.Parse(dobLayout, req.DOB)
	if err != nil {
		return models.UserResponse{}, err
	}

	user, err := s.repo.Create(ctx, req.Name, dob)
	if err != nil {
		return models.UserResponse{}, err
	}

	return models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		DOB:  user.Dob.Format(dobLayout),
	}, nil
}

func (s *userService) GetByID(ctx context.Context, id int32) (models.UserDetailResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.UserDetailResponse{}, ErrNotFound
		}
		return models.UserDetailResponse{}, err
	}

	return models.UserDetailResponse{
		ID:   user.ID,
		Name: user.Name,
		DOB:  user.Dob.Format(dobLayout),
		Age:  CalculateAge(user.Dob),
	}, nil
}

func (s *userService) Update(ctx context.Context, id int32, req models.UpdateUserRequest) (models.UserResponse, error) {
	// Verify the user exists.
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.UserResponse{}, ErrNotFound
		}
		return models.UserResponse{}, err
	}

	dob, err := time.Parse(dobLayout, req.DOB)
	if err != nil {
		return models.UserResponse{}, err
	}

	user, err := s.repo.Update(ctx, id, req.Name, dob)
	if err != nil {
		return models.UserResponse{}, err
	}

	return models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		DOB:  user.Dob.Format(dobLayout),
	}, nil
}

func (s *userService) Delete(ctx context.Context, id int32) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	return s.repo.Delete(ctx, id)
}

func (s *userService) List(ctx context.Context, page, pageSize int) (models.PaginatedUsersResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := int32((page - 1) * pageSize)
	users, err := s.repo.List(ctx, int32(pageSize), offset)
	if err != nil {
		return models.PaginatedUsersResponse{}, err
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return models.PaginatedUsersResponse{}, err
	}

	details := make([]models.UserDetailResponse, 0, len(users))
	for _, u := range users {
		details = append(details, models.UserDetailResponse{
			ID:   u.ID,
			Name: u.Name,
			DOB:  u.Dob.Format(dobLayout),
			Age:  CalculateAge(u.Dob),
		})
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	return models.PaginatedUsersResponse{
		Data:       details,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}
