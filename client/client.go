package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Client struct {
	backend string
	User    UserData
}

// NewClient create new API client.
func NewClient(backend string) *Client {
	var c Client
	c.backend = backend
	return &c
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

	//log.Println("json:", string(buf))

	var apiURL = "https://api.appnexus.com/auth"
	if c.backend == "simulator" {
		apiURL = "https://xandrtools.com/xandrsim/auth"
	}

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
		return fmt.Errorf("unmarshal: %w", err)
	}

	if respAuth.Response.Status != "OK" {
		return fmt.Errorf("error response [%s]", string(buf))
	}

	c.User.TokenData.Token = respAuth.Response.Token
	c.User.TokenData.ExpirationTime = time.Now().Add(time.Hour * 2)

	log.Println("auth request completed")

	return nil
}

// GetBatchSegmentJobs returns array of batch segment upload jobs
func (c *Client) GetBatchSegmentJobs(memberID int32) ([]BatchSegmentUploadJob, error) {
	if c.User.TokenData.Token == "" {
		return nil, fmt.Errorf("token is empty")
	}

	var apiURL = "https://api.appnexus.com/batch-segment"
	if c.backend == "simulator" {
		apiURL = "https://xandrtools.com/xandrsim/batch-segment"
	}

	log.Println("request:", apiURL)

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

	var bsResponse BatchSegmentResponse

	if err := json.Unmarshal(buf, &bsResponse); err != nil {
		log.Println("buf:", string(buf))
		return nil, err
	}

	log.Println("jobs:", len(bsResponse.Response.BatchSegmentUploadJob))

	return bsResponse.Response.BatchSegmentUploadJob, nil
}
