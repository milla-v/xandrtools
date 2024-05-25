package simulator

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestBatchSegment(t *testing.T) {

	var user UserData

	user.Username = "user1"
	user.TokenData.Token = "12345"
	user.TokenData.ExpirationTime = time.Now().Add(time.Second * 30)
	User.Store(user.TokenData.Token, user)

	//time.Sleep(8 * time.Second)

	testServer := httptest.NewServer(http.HandlerFunc(HandleBatchSegment))
	defer testServer.Close()

	req, err := http.NewRequest(http.MethodGet, testServer.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Authorization", "Bearer 12345")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	t.Log(string(buf))
}
