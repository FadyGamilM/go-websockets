package repository

import (
	"context"
	"database/sql"

	"github.com/FadyGamilM/go-websockets/internal/core"
	"github.com/FadyGamilM/go-websockets/internal/database/postgres"
	"github.com/FadyGamilM/go-websockets/internal/models"
)

type userRepo struct {
	pg *postgres.PG
}

func New(pg postgres.DBTX) core.UserRepository {
	return &userRepo{
		pg: &postgres.PG{
			DB: pg,
		},
	}
}

const (
	CREATE_QUERY = `
		INSERT INTO users 
		(username, email, password)
		VALUES
		($1, $2, $3)
		RETURNING id, username, email, password
	`

	GET_BY_ID_QUERY = `
		SELECT id, username, email, password
		FROM users
		WHERE id = $1
	`

	GET_BY_USERNAME_QUERY = `
		SELECT id, username, email, password
		FROM users
		WHERE username = $1
	`
)

func (ur *userRepo) Create(ctx context.Context, user *models.User) (*models.User, *core.DbError) {
	u := new(models.User)
	err := ur.pg.DB.QueryRowContext(ctx, CREATE_QUERY, user.Username, user.Email, user.Password).Scan(&u.ID, &u.Username, &u.Email, &u.Password)
	if err != nil {
		dbErr := core.New_DB_ErrorInsertingUser(user.Email)
		return nil, dbErr
	}

	return u, nil
}

func (ur *userRepo) GetByID(ctx context.Context, id int64) (*models.User, *core.DbError) {
	u := new(models.User)
	err := ur.pg.DB.QueryRowContext(ctx, GET_BY_ID_QUERY, id).Scan(&u.ID, &u.Username, &u.Email, &u.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			dbErr := core.New_DB_ErrorNonExistingUser(id)
			return nil, dbErr
		}
		dbErr := core.New_DB_ErrorFetchingUser(id)
		return nil, dbErr
	}
	return u, nil
}

func (ur *userRepo) GetByUsername(ctx context.Context, username string) (*models.User, *core.DbError) {
	u := new(models.User)
	err := ur.pg.DB.QueryRowContext(ctx, GET_BY_USERNAME_QUERY, username).Scan(&u.ID, &u.Username, &u.Email, &u.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			dbErr := core.New_DB_ErrorNonExistingUserWithUsername(username)
			return nil, dbErr
		}
		dbErr := core.New_DB_ErrorFetchingUserWithUsername(username)
		return nil, dbErr
	}
	return u, nil
}
