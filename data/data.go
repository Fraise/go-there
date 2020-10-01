package data

type User struct {
	Username string `db:"username" json:"username"`
	IsAdmin  bool   `db:"is_admin" json:"is_admin"`
	Password
	ApiKey
}

type Password struct {
	PasswordSalt []byte `db:"password_salt" json:"password_salt"`
	PasswordHash []byte `db:"password_hash" json:"password_hash"`
}

type ApiKey struct {
	ApiKeySalt []byte `db:"api_key_salt" json:"api_key_salt"`
	ApiKeyHash []byte `db:"api_key_hash" json:"api_key_hash"`
}
