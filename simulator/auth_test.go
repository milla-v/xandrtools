package simulator

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"xandrtools/client"
)

func TestAuthManyUsers(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(HandleAuthentication))
	defer testServer.Close()

	manyAuth := make([]client.AuthRequest, 5)

	for i := 0; i < len(manyAuth); i++ {
		manyAuth[i].Auth.Username = "user" + strconv.Itoa(i)
		manyAuth[i].Auth.Password = "password" + strconv.Itoa(i)
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

	var auth client.AuthRequest

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

func TestGenerateToken(t *testing.T) {
	token, err := generateToken(16)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("token: ", token)
}
