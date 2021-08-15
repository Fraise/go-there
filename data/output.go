package data

// ApiKeyResponse should be returned when creating a user or regenerating an API key.
type ApiKeyResponse struct {
	ApiKey string `json:"api_key,omitempty"`
}

// JwtResponse should be returned when querying the auth endpoint.
type JwtResponse struct {
	Jwt string `json:"jwt,omitempty"`
}

// ErrorResponse should be returned to the user when additional context is needed when an error occurs.
type ErrorResponse struct {
	Error string `json:"error"`
}
