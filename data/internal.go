package data

// User contains all the information representing an user internally. It should NOT be used to marshal/unmarshal
// incoming or outgoing data.
type User struct {
	Id           int    `db:"id" json:"id"`
	Username     string `db:"username" json:"username"`
	IsAdmin      bool   `db:"is_admin" json:"is_admin"`
	PasswordHash []byte `db:"password_hash" json:"password_hash,omitempty"`
	ApiKeyHash   []byte `db:"api_key_hash" json:"api_key_hash,omitempty"`
}

// Path contains the information representing a redirection target internally.
type Path struct {
	Path   string `db:"path" json:"path" binding:"required"`
	Target string `db:"target" json:"target" binding:"required"`
	UserId int    `db:"user_id"`
}

// LogInfo represents the data logged when a user makes a request.
type LogInfo struct {
	Method   string `json:"method"`
	Endpoint string `json:"endpoint"`
	User     string `json:"user"`
	Ip       string `json:"ip"`
	HttpCode int    `json:"http_code"`
}
