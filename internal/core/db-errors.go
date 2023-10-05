package core

import "fmt"

type DB_ERROR string

const (
	ERROR_INSERTING_USER DB_ERROR = "error inserting new user"
	ERROR_FETCHING_USER  DB_ERROR = "error retrieving a user"
	ERROR_NOT_FOUND_USER DB_ERROR = "error non existing user"
)

type AppError struct {
	Msg  string
	Type DB_ERROR
}

func (ae *AppError) Error() string {
	return ae.Msg
}

func New_ErrorInsertingUser(email string) *AppError {
	return &AppError{
		Msg:  fmt.Sprintf("error inserting a new user with email = %v", email),
		Type: ERROR_INSERTING_USER,
	}
}

func New_ErrorFetchingUser(id int64) *AppError {
	return &AppError{
		Msg:  fmt.Sprintf("error selecting a user with id = %v", id),
		Type: ERROR_FETCHING_USER,
	}
}

func New_ErrorNonExistingUser(id int64) *AppError {
	return &AppError{
		Msg:  fmt.Sprintf("error not found user with id = %v", id),
		Type: ERROR_NOT_FOUND_USER,
	}
}
