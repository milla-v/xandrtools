package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

	log.Println("json:", string(buf))

	var apiURL = "https://api.appnexus.com/auth"
	if c.backend == "simulator" {
		apiURL = "http://127.0.0.1:9970/xandrsim/auth"
	}

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

	log.Println("token:", respAuth.Response.Token)
	c.User.TokenData.Token = respAuth.Response.Token
	c.User.TokenData.ExpirationTime = time.Now().Add(time.Hour * 2)

	return nil
}

func (c *Client) GetBatchSegmentJobs(user UserData) error {
	if user.Username == "" {
		return fmt.Errorf("username is empty")
	}
	if time.Time.IsZero(user.TokenData.ExpirationTime) == true {
		return fmt.Errorf("expiration time is empty")
	}
	if user.TokenData.Token == "" {
		return fmt.Errorf("token is empty")
	}
	if user.TokenData.MemberId == "" {
		return fmt.Errorf("member_id is empty")
	}
	c.User.Username = user.Username
	c.User.TokenData.Token = user.TokenData.Token
	c.User.TokenData.ExpirationTime = user.TokenData.ExpirationTime
	c.User.TokenData.MemberId = user.TokenData.MemberId

	var apiURL = "https://api.appnexus.com/batch-segment"
	if c.backend == "simulator" {
		apiURL = "http://127.0.0.1:9970/xandrsim/batch-segment"
	}
	log.Println("apiURL", apiURL)

	return nil
}
