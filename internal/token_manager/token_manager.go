package tokenmanager

import "banner_service/internal/domains"

type Token struct {
	Valid bool
	role  domains.Role
}

func (t *Token) GetRole() domains.Role {
	return t.role
}

type TokenManager struct {
	userToken  string
	adminToken string
}

func New(userToken, adminToken string) *TokenManager {
	return &TokenManager{userToken, adminToken}
}

func (tm *TokenManager) Parse(inputToken string) Token {
	t := Token{}
	if inputToken == tm.userToken {
		t.role = domains.UserRole
		t.Valid = true
		return t
	}

	if inputToken == tm.adminToken {
		t.role = domains.AdminRole
		t.Valid = true
		return t
	}

	return t
}
