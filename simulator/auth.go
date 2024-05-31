package simulator

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"xandrtools/client"
)

var User sync.Map

func HandleAuthentication(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "error", http.StatusMethodNotAllowed)
		return
	}

	buf, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	var auth client.AuthRequest
	var user client.UserData

	if err := json.Unmarshal(buf, &auth); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("login user", auth.Auth.Username, ":", auth.Auth.Password)

	var authResp client.AuthResponse
	if auth.Auth.Username != "" || auth.Auth.Password != "" {
		authResp.Response.Status = "OK"
		authResp.Response.Token, err = generateToken(16)
		if err != nil {
			log.Println("generate token err: ", err)
			return
		}
	}

	//fill UserData struct
	user.TokenData.Token = authResp.Response.Token
	//	log.Println("RANDOM INT", rand.IntN(10))

	user.TokenData.ExpirationTime = time.Now().Add(time.Hour * 2) //token expiration time - 2 hours

	User.Store(user.TokenData.Token, user)

	buf, err = json.MarshalIndent(authResp, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := fmt.Fprintln(w, string(buf)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//fmt.Printf("json data: %s\n", buf)
}
