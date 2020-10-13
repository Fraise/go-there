package data

import "errors"

// Datasource errors
var (
	ErrSql             = errors.New("sql: error")
	ErrSqlNoRow        = errors.New("sql: no row in result set")
	ErrSqlDuplicateRow = errors.New("sql: duplicate row")
)

// Auth errors
var (
	ErrInvalidKey = errors.New("auth: invalid api key")
)
