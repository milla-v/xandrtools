package simulator

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
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

	var auth AuthRequest

	if err := json.Unmarshal(buf, &auth); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("login user", auth.Auth.Username, ":", auth.Auth.Password)

	var authResp AuthResponse

	authResp.Response.Status = "OK"
	authResp.Response.Token = "12345"

	buf, err = json.MarshalIndent(authResp, "\t", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := io.Copy(w, bytes.NewReader(buf)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
