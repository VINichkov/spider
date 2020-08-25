package models

import "errors"

var (
	ErrNotFound = errors.New("sql: no rows in result set")
)

func EqualErrors(err1 error, err2 error) bool{
	if err1 != nil && err2 != nil {
		return err1.Error() == err2.Error()
	}
	return false
}