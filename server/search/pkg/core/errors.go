package core

import "errors"

var (
	ErrCategoryAlreadyExists = errors.New("category already exists")
	ErrCategoryDoesNotExists = errors.New("category does not exists")
)
