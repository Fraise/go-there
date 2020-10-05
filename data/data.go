package data

// User contains all the informations representing an user internally. It should NOT be used to marshal/unmarshal
// incoming or outgoing data.
type User struct {
	Username     string `db:"username" json:"username"`
	IsAdmin      bool   `db:"is_admin" json:"is_admin"`
	PasswordHash []byte `db:"password_hash" json:"password_hash"`
	ApiKeySalt   []byte `db:"api_key_salt" json:"api_key_salt"`
	ApiKeyHash   []byte `db:"api_key_hash" json:"api_key_hash"`
}

// Login represents the information given by a user to authenticate. It should be used to unmarshal incoming
// authentication data.
type Login struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
	ApiKey   string `form:"api_key" json:"api_key"`
}

// CreateUser represents the information given by a user to create another user. It should be used to unmarshal incoming
// creation data.
type CreateUser struct {
	CreateUser     string `json:"create_user" binding:"required"`
	CreatePassword string `json:"create_password" binding:"required"`
}

// ApiKeyResponse should be returned when creating a user or regenerating an API key.
type ApiKeyResponse struct {
	ApiKey string `json:"api_key"`
}

// ErrorResponse should be returned to the user when additional context is needed when an error occurs.
type ErrorResponse struct {
	Error string `json:"error"`
}
