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

type BatchSegment struct {
	Username string `json: "username"`
	Segment  struct {
		MemberId string `json: "memberid"`
		Status   bool   `json: "status`
	} `json:"segment"`
}

type UserData struct {
	Username  string `json: "username"`
	TokenData struct {
		Token          string    `json: "token"`
		ExpirationTime time.Time `json: "expirationTime"`
		MemberId       string    `json: "memberid"`
	} `json: "tokendata"`
}
