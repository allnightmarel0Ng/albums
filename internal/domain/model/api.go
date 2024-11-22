package model

type AuthorizationResponse struct {
	Code  int    `json:"code"`
	Jwt   string `json:"jwt,omitempty"`
	Error string `json:"error,omitempty"`
}
