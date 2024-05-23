package simulator

import (
	"log"
	"net/http"
	"strings"
	"time"
)

func HandleBatchSegment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "GET", http.StatusMethodNotAllowed)
		return
	}

	for k, v := range r.Header {
		log.Printf("header: %s=%v", k, v)
	}

	tokenHeader := r.Header.Get("Authorization")
	token := strings.TrimPrefix(tokenHeader, "Bearer ")

	log.Println("token:", token)

	user, ok := User.Load(token)
	if !ok {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	log.Printf("user: %+v", user)
	log.Println("token expiration time: ", user)

	u := user.(UserData)
	log.Println("u.TokenData. ExpirationTime: ", u.TokenData.ExpirationTime)

	// TODO: check expiration time
	if u.TokenData.ExpirationTime.IsZero() == true {
		http.Error(w, "invalid expiration time: ", http.StatusUnauthorized)
		return
	} else {
		if time.Now().Before(u.TokenData.ExpirationTime) == false {
			http.Error(w, "invalid expiration time: ", http.StatusUnauthorized)
			return
		}
	}

	var resp BatchSegmentResponse
	_ = resp
}
