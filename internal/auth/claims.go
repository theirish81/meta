package auth

import "github.com/golang-jwt/jwt/v5"

type MetaClaims struct {
	jwt.RegisteredClaims
	Email       string `json:"email"`
	Nonce       string `json:"nonce"`
	Permissions string `json:"permission"`
}

func (m MetaClaims) CanRead() bool {
	return m.Permissions == "admin" || m.Permissions == "write" || m.Permissions == "read"
}

func (m MetaClaims) CanWrite() bool {
	return m.Permissions == "admin" || m.Permissions == "write"
}

func (m MetaClaims) CanAdmin() bool {
	return m.Permissions == "admin"
}
