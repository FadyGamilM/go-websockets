package core

import "fmt"

type DB_ERROR string

const (
	// both are internal server error at the end
	ERROR_INSERTING_USER DB_ERROR = "error inserting new user"
	ERROR_FETCHING_USER  DB_ERROR = "error retrieving a user"
	// bad request error (client error)
	ERROR_NOT_FOUND_USER DB_ERROR = "error non existing user"
)

type DbError struct {
	Msg  string
	Type DB_ERROR
}

func (ae *DbError) Error() string {
	return ae.Msg
}

func New_DB_ErrorInsertingUser(email string) *DbError {
	return &DbError{
		Msg:  fmt.Sprintf("error inserting a new user with email = %v", email),
		Type: ERROR_INSERTING_USER,
	}
}

func New_DB_ErrorFetchingUser(id int64) *DbError {
	return &DbError{
		Msg:  fmt.Sprintf("error selecting a user with id = %v", id),
		Type: ERROR_FETCHING_USER,
	}
}

func New_DB_ErrorFetchingUserWithUsername(username string) *DbError {
	return &DbError{
		Msg:  fmt.Sprintf("error selecting a user with username = %v", username),
		Type: ERROR_FETCHING_USER,
	}
}

func New_DB_ErrorNonExistingUser(id int64) *DbError {
	return &DbError{
		Msg:  fmt.Sprintf("error not found user with id = %v", id),
		Type: ERROR_NOT_FOUND_USER,
	}
}

func New_DB_ErrorNonExistingUserWithUsername(username string) *DbError {
	return &DbError{
		Msg:  fmt.Sprintf("error not found user with username = %v", username),
		Type: ERROR_NOT_FOUND_USER,
	}
}
