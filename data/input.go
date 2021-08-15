package data

// UserInfo contains the name and redirections created by an user.
type UserInfo struct {
	Username string     `db:"username" json:"username"`
	IsAdmin  bool       `db:"is_admin" json:"is_admin"`
	Paths    []PathInfo `json:"paths,omitempty"`
}

// PathInfo contains the pair Path/Target.
type PathInfo struct {
	Path   string `db:"path" json:"path,omitempty" binding:"required"`
	Target string `db:"target" json:"target,omitempty" binding:"required"`
}

// CreatePath represents the data sent by the user to add a new redirection path.
type CreatePath struct {
	Path   string `json:"path" binding:"required"`
	Target string `json:"target" binding:"required"`
}

// DeletePath represents the data sent by the user to delete an existing redirection path.
type DeletePath struct {
	Path string `json:"path" binding:"required"`
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
