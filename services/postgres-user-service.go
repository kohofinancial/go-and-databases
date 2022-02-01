package services

import (
	"context"
	"database/sql"
	"github.com/gofrs/uuid"
	"time"
)

type PostgresUserService struct {
	db *sql.DB
}

func NewPostgresUserService(db *sql.DB) *PostgresUserService {
	return &PostgresUserService{db: db}
}

func (p *PostgresUserService) Get(id string) (*User, error) {
	query := `
        SELECT id, name, occupation, created_at, updated_at
        FROM users
        WHERE id = $1`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := p.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Occupation,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *PostgresUserService) Delete(id string) error {
	query := `
        DELETE FROM users
        WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result, err := p.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (p *PostgresUserService) DeleteAll() error {
	//goland:noinspection SqlWithoutWhere
	query := `DELETE FROM users`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result, err := p.db.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (p *PostgresUserService) Update(user *User) error {
	query := `
        UPDATE users 
        SET name = $2, occupation = $3
        WHERE id = $1
        RETURNING updated_at`

	args := []interface{}{
		user.ID,
		user.Name,
		user.Occupation,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := p.db.QueryRowContext(ctx, query, args...).Scan(&user.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresUserService) Create(user *User) error {
	query := `
        INSERT INTO users (id, name, occupation) 
        VALUES ($1, $2, $3)
        RETURNING id, created_at, updated_at`

	newV4, err := uuid.NewV4()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	return p.db.QueryRowContext(ctx, query, newV4.String(), user.Name, user.Occupation).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}
