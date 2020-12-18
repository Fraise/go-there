package data

import "errors"

// Datasource errors
var (
	ErrSql             = errors.New("sql: error")
	ErrSqlNoRow        = errors.New("sql: no row in result set")
	ErrSqlDuplicateRow = errors.New("sql: duplicate row")
	ErrRedis           = errors.New("redis: error")
)

// Auth errors
var (
	ErrInvalidKey = errors.New("auth: invalid api key")
)

// Init errors
var (
	ErrInit = errors.New("init: failed")
)
