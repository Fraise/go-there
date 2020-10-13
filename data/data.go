package data

// User contains all the informations representing an user internally. It should NOT be used to marshal/unmarshal
// incoming or outgoing data.
type User struct {
	Username     string `db:"username" json:"username"`
	IsAdmin      bool   `db:"is_admin" json:"is_admin"`
	PasswordHash []byte `db:"password_hash" json:"password_hash,omitempty"`
	ApiKeySalt   []byte `db:"api_key_salt" json:"api_key_salt,omitempty"`
	ApiKeyHash   []byte `db:"api_key_hash" json:"api_key_hash,omitempty"`
}

// Path contains the informations representing a redirection target internally.
type Path struct {
	Path   string `db:"path" json:"path" binding:"required"`
	Target string `db:"target" json:"target" binding:"required"`
	User   string `db:"user"`
}

// CreatePath represents the data sent by the user to add a new redirection path.
type CreatePath struct {
	Path   string `json:"path" binding:"required"`
	Target string `json:"target" binding:"required"`
}

// DeletePath represents the data sent by the user to delete an existing redirection path.
type DeletePath struct {
	Path string `json:"path"`
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

// PatchUser represents the input used to change a user password or request a new API key.
type PatchUser struct {
	PatchPassword string `json:"new_password"`
	PatchApiKey   bool   `json:"new_api_key"`
}

// ApiKeyResponse should be returned when creating a user or regenerating an API key.
type ApiKeyResponse struct {
	ApiKey string `json:"api_key,omitempty"`
}

// ErrorResponse should be returned to the user when additional context is needed when an error occurs.
type ErrorResponse struct {
	Error string `json:"error"`
}
