package data

type User struct {
	Username     string `db:"username" json:"username"`
	IsAdmin      bool   `db:"is_admin" json:"is_admin"`
	PasswordHash []byte `db:"password_hash" json:"password_hash"`
	ApiKeySalt   []byte `db:"api_key_salt" json:"api_key_salt"`
	ApiKeyHash   []byte `db:"api_key_hash" json:"api_key_hash"`
}

type Login struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
	ApiKey   string `form:"api_key" json:"api_key"`
}

type CreateUser struct {
	CreateUser     string `json:"create_user" binding:"required"`
	CreatePassword string `json:"create_password" binding:"required"`
}

type CreateUserResponse struct {
	ApiKey string `json:"api_key"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
