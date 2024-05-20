package simulator

import (
	"time"
)

type AuthRequest struct {
	Auth struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"auth"`
}

type AuthResponse struct {
	Response struct {
		Status string `json:"status"`
		Token  string `json:"token"`
	} `json:"response"`
}

type UserData struct {
	Username  string
	TokenData struct {
		Token          string
		ExpirationTime time.Time
	}
}
