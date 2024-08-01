package simulator

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"xandrtools/client"
)

func TestBatchSegment(t *testing.T) {
	createTestUser()

	//time.Sleep(8 * time.Second)

	testServer := httptest.NewServer(http.HandlerFunc(HandleBatchSegment))
	defer testServer.Close()
	os.Setenv("TEST_BACKEND", testServer.URL)

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

func TestClientNoAuthError(t *testing.T) {
	createTestUser()

	testServer := httptest.NewServer(http.HandlerFunc(HandleBatchSegment))
	defer testServer.Close()
	os.Setenv("TEST_BACKEND", testServer.URL)

	cli := client.NewClient(testServer.URL)
	cli.User.TokenData.Token = "123"
	_, err := cli.GetBatchSegmentJobs(111)
	if err == nil {
		t.Fatal("should be NOAUTH error")
	}
	if !strings.Contains(err.Error(), "NOAUTH:Authentication failed - not logged in") {
		t.Fatal("should be auth error but returned: ", err.Error())
	}
}

func TestClientNoMemberIdError(t *testing.T) {
	createTestUser()

	testServer := httptest.NewServer(http.HandlerFunc(HandleBatchSegment))
	defer testServer.Close()
	os.Setenv("TEST_BACKEND", testServer.URL)

	cli := client.NewClient("")
	cli.User.TokenData.Token = "12345"
	_, err := cli.GetBatchSegmentJobs(0)
	if err == nil {
		t.Fatal("should be SYNTAX error")
	}
	if !strings.Contains(err.Error(), "SYNTAX:no member_id provided") {
		t.Fatal("should syntax error but returned: ", err.Error())
	}
}

func TestClientGetJobsSuccess(t *testing.T) {
	createTestUser()

	testServer := httptest.NewServer(http.HandlerFunc(HandleBatchSegment))
	defer testServer.Close()
	os.Setenv("TEST_BACKEND", testServer.URL)

	cli := client.NewClient("")
	cli.User.TokenData.Token = "12345"
	jobs, err := cli.GetBatchSegmentJobs(111)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("jobs:", len(jobs))
}

func createTestUser() {
	var user client.UserData

	user.Username = "user1"
	user.TokenData.Token = "12345"
	user.TokenData.ExpirationTime = time.Now().Add(time.Second * 30)
	User.Store(user.TokenData.Token, user)
}
