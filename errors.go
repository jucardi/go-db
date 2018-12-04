package dbx

import "fmt"

const (
	ErrNotFound = ErrType(iota >> 1)
	ErrFileAccess
	ErrDbAccess
	ErrDbOperation
	ErrMigrationFailed
)

type ErrType int

type DbError struct {
	Code    ErrType
	Message string
}

func (err *DbError) Error() string {
	return fmt.Sprintf("%s", err.Message)
}

func (err *DbError) String() string {
	return err.Error()
}

func (err *DbError) Is(t ErrType) bool {
	return err.Code|t == t
}

func (err *DbError) IsNotFound() bool {
	return err.Code == ErrNotFound
}
