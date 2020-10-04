package data

import "errors"

// Datasource errors
var ErrSqlNoRow = errors.New("sql: no rows in result set")

// Auth errors
var ErrInvalidKey = errors.New("invalid api key")
