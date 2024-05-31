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

func (c *Client) GetBatchSegmentJobs() error {
	return nil
}
