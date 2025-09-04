package repository

import "errors"

var (
	ErrNotFound   = errors.New("repository: not found")
	ErrConflict   = errors.New("repository: conflict")
	ErrBadRequest = errors.New("repository: bad request")
)
