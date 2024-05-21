package simulator

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"sync"
	"time"
)

var UserToken sync.Map
var User sync.Map
var UserRequest sync.Map

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
	var user UserData

	if err := json.Unmarshal(buf, &auth); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("login user", auth.Auth.Username, ":", auth.Auth.Password)

	var authResp AuthResponse
	if auth.Auth.Username != "" || auth.Auth.Password != "" {
		authResp.Response.Status = "OK"
		authResp.Response.Token, err = generateToken()
		if err != nil {
			log.Println("generate token err: ", err)
			return
		}
	}

	//fill UserData struct
	user.TokenData.Token = authResp.Response.Token
	log.Println("RANDOM INT", rand.IntN(10))

	user.TokenData.ExpirationTime = time.Now().Add(time.Hour * 2) //token expiration time - 2 hours

	User.Store(auth.Auth.Username, user.TokenData)
	if userValue, ok := User.Load(auth.Auth.Username); ok {
		log.Printf("Key %s - Value %d\n", auth.Auth.Username, userValue)
	}
	UserToken.Store(auth.Auth.Username, authResp.Response.Token)
	/*
		if value, ok := UserToken.Load(auth.Auth.Username); ok {
			log.Printf("Key %s - Value %d\n", auth.Auth.Username, value)
		}
	*/

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

func HandleBatchSegmentRequest(w http.ResponseWriter, r *http.Request) {
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

	log.Println("r.Method = ", r.Method)

	var user UserData

	if err := json.Unmarshal(buf, &user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var batchSeg BatchSegment
	//check if userID exists
	if user.Username == "" || user.TokenData.MemberId == "" || user.TokenData.Token == "" {
		log.Println("wrong data")
		return
	}

	t := time.Now()
	log.Println("time now: ", t.Second(), "expiration time: ", user.TokenData.ExpirationTime.Second())
	batchSeg.Segment.Status = t.Before(user.TokenData.ExpirationTime)
	batchSeg.Segment.MemberId = user.TokenData.MemberId
	log.Println("username: ", user.Username, "batchSeg.Segment.MemberId = ", batchSeg.Segment.MemberId, "STATUS: ", batchSeg.Segment.Status)
	UserRequest.Store(user.Username, batchSeg)
	buf, err = json.MarshalIndent(batchSeg, "\t", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := io.Copy(w, bytes.NewReader(buf)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
