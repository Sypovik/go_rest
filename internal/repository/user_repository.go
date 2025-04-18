package repository

import (
	"context"
	"database/sql"
	"go-restapi/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetById(ctx context.Context, id int) (*model.User, error)
	GetByName(ctx context.Context, id int) (*model.User, error)
	GetAll(ctx context.Context) ([]model.User, error)
	// Update(ctx context.Context, user *model.User) error
	// Delete(ctx context.Context, id int) error
	Count(ctx context.Context) (int, error)
}

type userRepository struct {
	db *sql.DB
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	query := `
	INSERT INTO user (name, password)
	VALUES ($1? $2)
	RETURNING ID
	`

	return r.db.QueryRowContext(ctx, query, user.Name, user.Password).Scan(&user.ID)
}
func (r *userRepository) GetById(ctx context.Context, id int) (*model.User, error) {
	query := `
	SELECT * FROM User
	WHERE ID = $1
	`

	var user model.User
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name, &user.Password); err != nil {
		return nil, err

	}
	return &user, nil
}

func (r *userRepository) GetByName(ctx context.Context, name string) (*model.User, error) {
	query := `
	SELECT * FROM User
	WHERE Name = $1
	`

	var user model.User
	if err := r.db.QueryRowContext(ctx, query, name).Scan(&user.ID, &user.Name, &user.Password); err != nil {
		return nil, err

	}
	return &user, nil
}

func (r *userRepository) GetAll(ctx context.Context) ([]model.User, error) {
	query := `
	SELECT * FROM User
	ORDER BY ID
	`
	// var users[] model.User

	var users []model.User
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Password); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *userRepository) Count(ctx context.Context) (int, error) {
	var count int
	query := `
	SELECT COUNT(*)
	FROM User
	`
	if err := r.db.QueryRowContext(ctx, query).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil

}
