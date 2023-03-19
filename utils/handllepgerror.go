package utils

import (
	"errors"
	"net/http"

	"github.com/uptrace/bun/driver/pgdriver"
)

func HandlePostgresError(err error) (error, bool) {
	if err, ok := err.(pgdriver.Error); ok {
		switch err.Field('C') {
		case "23505":
			// Unique constraint violation
			return CustomAPIErr{Code: http.StatusBadRequest, Err: errors.New("value already exists in table " + err.Field('t') + " for column " + err.Field('n'))}, true
		case "42P01":
			// Table not found
			return CustomAPIErr{Code: http.StatusBadRequest, Err: errors.New("table " + err.Field('t') + "doesn't exists.")}, true
		default:
			// Unhandled error code
			return CustomAPIErr{Code: http.StatusInternalServerError, Err: err}, true
		}
	}
	return err, false
}
