package simulator

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
