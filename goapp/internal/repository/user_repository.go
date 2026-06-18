package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/yash/user-dob-api/db/sqlc"
)

// UserRepository defines the data-access contract.
type UserRepository interface {
	Create(ctx context.Context, name string, dob time.Time) (sqlc.User, error)
	GetByID(ctx context.Context, id int32) (sqlc.User, error)
	Update(ctx context.Context, id int32, name string, dob time.Time) (sqlc.User, error)
	Delete(ctx context.Context, id int32) error
	List(ctx context.Context, limit, offset int32) ([]sqlc.User, error)
	Count(ctx context.Context) (int64, error)
}

type userRepository struct {
	q *sqlc.Queries
}

// New returns a new UserRepository backed by db.
func New(db *sql.DB) UserRepository {
	return &userRepository{q: sqlc.New(db)}
}

func (r *userRepository) Create(ctx context.Context, name string, dob time.Time) (sqlc.User, error) {
	return r.q.CreateUser(ctx, sqlc.CreateUserParams{
		Name: name,
		Dob:  dob,
	})
}

func (r *userRepository) GetByID(ctx context.Context, id int32) (sqlc.User, error) {
	return r.q.GetUser(ctx, id)
}

func (r *userRepository) Update(ctx context.Context, id int32, name string, dob time.Time) (sqlc.User, error) {
	return r.q.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:   id,
		Name: name,
		Dob:  dob,
	})
}

func (r *userRepository) Delete(ctx context.Context, id int32) error {
	return r.q.DeleteUser(ctx, id)
}

func (r *userRepository) List(ctx context.Context, limit, offset int32) ([]sqlc.User, error) {
	return r.q.ListUsers(ctx, sqlc.ListUsersParams{
		Limit:  limit,
		Offset: offset,
	})
}

func (r *userRepository) Count(ctx context.Context) (int64, error) {
	return r.q.CountUsers(ctx)
}
