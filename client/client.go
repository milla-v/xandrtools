package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Client struct {
	backend string
	User    UserData
	log     *log.Logger
}

// NewClient create new API client.
func NewClient(backend string) *Client {
	var c Client
	c.backend = backend

	debugAddr := os.Getenv("DEBUG_ADDR")
	if debugAddr != "" {
		c.log = log.New(os.Stderr, "client: ", log.Lmsgprefix|log.Lshortfile|log.LstdFlags)
	} else {
		c.log = log.New(io.Discard, "", 0)
	}

	return &c
}

func (c *Client) getAPIURL() string {
	var apiURL = "https://api.appnexus.com/"
	if c.backend == "simulator" {
		apiURL = "https://xandrtools.com/xandrsim/"
	}
	testUrl := os.Getenv("TEST_BACKEND")
	if testUrl != "" {
		apiURL = testUrl + "/"
	}
	return apiURL
}

// Login returns user token.
func (c *Client) Login(username, password string) error {
	var auth AuthRequest

	auth.Auth.Username = username
	auth.Auth.Password = password

	if auth.Auth.Username == "" {
		return fmt.Errorf("username is empty")
	}

	if auth.Auth.Password == "" {
		return fmt.Errorf("password is empty")
	}

	buf, err := json.MarshalIndent(auth, "\t", "\t")
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	apiURL, err := getApiURL(c.backend)
	if err != nil {
		c.log.Println("get api url err:", err)
		return nil
	}
	log.Println("I AM HERE", " apiUrl: ", apiURL)
	apiURL += "auth"

	log.Println("request:", apiURL, "user:", username)

	resp, err := http.Post(apiURL, "application/json", bytes.NewReader(buf))
	if err != nil {
		return fmt.Errorf("post: %w", err)
	}

	buf, err = io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("readall: %w", err)
	}
	defer resp.Body.Close()

	var respAuth AuthResponse

	if err := json.Unmarshal(buf, &respAuth); err != nil {
		return fmt.Errorf("unmarshal: %w.\nresp: %s", err, string(buf))
	}

	if respAuth.Response.Status != "OK" {
		return fmt.Errorf("error response [%s]", string(buf))
	}

	c.User.TokenData.Token = respAuth.Response.Token
	c.User.TokenData.ExpirationTime = time.Now().UTC().Add(time.Hour * 2)
	c.log.Println("auth request completed")

	return nil
}

// GetBatchSegmentJobs returns array of batch segment upload jobs
func (c *Client) GetBatchSegmentJobs(memberID int32) ([]BatchSegmentUploadJob, error) {
	if c.User.TokenData.Token == "" {
		return nil, fmt.Errorf("token is empty")
	}
	var apiURL string
	var bsResponse BatchSegmentResponse
	c.log.Println("BACKEND Get Job: ", c.backend)

	apiURL, err := getApiURL(c.backend)
	if err != nil {
		c.log.Println("get api url err:", err)
		return bsResponse.Response.BatchSegmentUploadJob, nil
	}

	apiURL += "batch-segment?member_id=" + strconv.Itoa(int(memberID))

	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+c.User.TokenData.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var errResponse ErrorResponse
	if err := json.Unmarshal(buf, &errResponse); err != nil {
		c.log.Println("buf:", string(buf))
		return nil, err
	}

	if errResponse.Response.Error != "" {
		return nil, fmt.Errorf("%s:%s", errResponse.Response.ErrorId, errResponse.Response.Error)
	}

	if err := json.Unmarshal(buf, &bsResponse); err != nil {
		c.log.Println("buf:", string(buf))
		return nil, err
	}

	//c.log.Println("jobs:", len(bsResponse.Response.BatchSegmentUploadJob))
	//c.log.Println("COMPLETE_time: ", bsResponse.Response.BatchSegmentUploadJob[0].CompletedTime)

	return bsResponse.Response.BatchSegmentUploadJob, nil
}
