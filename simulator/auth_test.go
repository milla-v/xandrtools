package simulator

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func TestUserRequest(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(HandleUserRequest))
	defer testServer.Close()

	//create 5 users
	userRequest := make([]UserData, 5)
	for i := 0; i < len(userRequest); i++ {
		userRequest[i].Username = "username" + strconv.Itoa(i)
		userRequest[i].TokenData.MemberId = strconv.Itoa(i+1) + strconv.Itoa(2) + strconv.Itoa(i+3)
		userRequest[i].TokenData.Token, _ = generateToken()
		userRequest[i].TokenData.ExpirationTime = time.Now().Add(time.Second * 2)
	}

	for i := 0; i < len(userRequest); i++ {
		buf, err := json.MarshalIndent(userRequest[i], "\t", "\t")
		if err != nil {
			t.Fatal(err)
		}
		resp, err := http.Post(testServer.URL, "application/json", bytes.NewReader(buf))
		if err != nil {
			t.Fatal(err)
		}

		buff, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		t.Logf("status:%s body: %s", resp.Status, string(buff))
		time.Sleep(5 * time.Second)
	}

}

func TestAuthManyUsers(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(HandleAuthentication))
	defer testServer.Close()

	manyAuth := make([]AuthRequest, 5)

	for i := 0; i < len(manyAuth); i++ {
		manyAuth[i].Auth.Username = "user" + strconv.Itoa(i)
		manyAuth[i].Auth.Password = "psssword" + strconv.Itoa(i)
	}

	for i := 0; i < len(manyAuth); i++ {
		t.Logf("user: %s | password: %s", manyAuth[i].Auth.Username, manyAuth[i].Auth.Password)
	}

	for i := 0; i < len(manyAuth); i++ {
		buf, err := json.MarshalIndent(manyAuth[i], "\t", "\t")
		if err != nil {
			t.Fatal(err)
		}
		resp, err := http.Post(testServer.URL, "application/json", bytes.NewReader(buf))
		if err != nil {
			t.Fatal(err)
		}

		buff, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		t.Logf("status:%s body: %s", resp.Status, string(buff))
	}
}

func TestAuthSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(HandleAuthentication))
	defer testServer.Close()

	var auth AuthRequest

	auth.Auth.Username = "user1"
	auth.Auth.Password = "pass1"

	buf, err := json.MarshalIndent(auth, "\t", "\t")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("json:", string(buf))

	resp, err := http.Post(testServer.URL, "application/json", bytes.NewReader(buf))
	if err != nil {
		t.Fatal(err)
	}

	buff, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	t.Logf("status:%s body: %s", resp.Status, string(buff))
}

func TestAuthError(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(HandleAuthentication))
	defer testServer.Close()

	resp, err := http.Get(testServer.URL)
	if err != nil {
		t.Fatal(err)
	}

	buff, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	t.Logf("status:%s body: %s", resp.Status, string(buff))

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatal(err)
	}
}

/*
	func TestAuthPostRequest(t *testing.T) {
		var user = Auth{
			Username: "authUser",
			Password: "authPassword",
		}
		err := authPostRequest(user)
		if err != nil {
			t.Fatal(err)
		}
	}

	func TestWriteAuthFile(t *testing.T) {
		var user = Auth{
			Username: "authuser",
			Password: "authpassword",
		}
		filename := "auth.json"
		var err error
		if err = writeAuthFile(user, filename); err != nil {
			t.Fatal(err)
		}
	}

	func TestReadAuthFile(t *testing.T) {
		var user Auth
		var err error
		fileName := "auth.json"
		user, err = readAuthFile(fileName)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("Unmarshaled user data: ", user)
	}
*/
func TestGenerateToken(t *testing.T) {
	token, err := generateToken()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("token: ", token)
}
