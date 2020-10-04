package data

type User struct {
	Username     string `db:"username" json:"username"`
	IsAdmin      bool   `db:"is_admin" json:"is_admin"`
	PasswordHash []byte `db:"password_hash" json:"password_hash"`
	ApiKeyHash   []byte `db:"api_key_hash" json:"api_key_hash"`
}

type Login struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
	ApiKey   string `form:"api_key" json:"api_key"`
}

type CreateUser struct {
	CreateUsername string `json:"create_username"`
	CreatePassword string `json:"create_password"`
}
