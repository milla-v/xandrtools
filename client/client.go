package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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

	var isDebug bool
	log.Println("BACKEND Login: ", c.backend)
	debugAddr := os.Getenv("DEBUG_ADDR")
	if debugAddr != "" {
		isDebug = true
	}
	log.Println("ADDR in Login: ", debugAddr)

	log.Println("c.backend: ", c.backend)
	var apiURL string
	switch {
	case c.backend == "simulator" && isDebug == true:
		apiURL = "https://" + debugAddr + "/xandrsim/auth"
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		log.Println("apiUPL: ", apiURL)
	case c.backend == "simulator" && isDebug == false:
		apiURL = "https://xandrtools.com/xandrsim/auth"
		log.Println("apiUPL: ", apiURL)
	case c.backend == "xandr" && isDebug == false:
		apiURL = "https://api.appnexus.com/auth"
		log.Println("apiUPL: ", apiURL)
	}

	/*
		var apiURL = "https://api.appnexus.com/auth"
		if c.backend != "" {
			apiURL = "https://xandrtools.com/xandrsim/auth"
		}
		if strings.HasPrefix(c.backend, "test:") {
			apiURL = strings.TrimPrefix(c.backend, "test:")
		}
	*/
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
	log.Println("auth request completed")

	return nil
}

// GetBatchSegmentJobs returns array of batch segment upload jobs
func (c *Client) GetBatchSegmentJobs(memberID int32) ([]BatchSegmentUploadJob, error) {
	if c.User.TokenData.Token == "" {
		return nil, fmt.Errorf("token is empty")
	}
	var apiURL string
	var isDebug bool
	log.Println("BACKEND Get Job: ", c.backend)
	debugAddr := os.Getenv("DEBUG_ADDR")
	if debugAddr != "" {
		isDebug = true
	}
	log.Println("ADDR in jobs: ", debugAddr)

	if strings.HasPrefix(c.backend, "test:") {
		apiURL = strings.TrimPrefix(c.backend, "test:")
		log.Println("URL: ", apiURL)
	}
	switch {
	case c.backend == "simulator" && isDebug == true:
		apiURL = "https://" + debugAddr + "/xandrsim/batch-segment"
		log.Println("apiUPL: ", apiURL)
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	case c.backend == "simulator" && isDebug == false:
		apiURL = "https://xandrtools.com/xandrsim/batch-segment"
		log.Println("apiUPL: ", apiURL)
	case c.backend == "xandr" && isDebug == false:
		apiURL = "https://api.appnexus.com/batch-segment"
		log.Println("apiUPL: ", apiURL)
	}
	/*
		if c.backend == "simulator" {
			apiURL = "https://xandrtools.com/xandrsim/batch-segment"
			log.Println("apiUPL: ", apiURL)
		}
		if c.backend == "xandr" {
			apiURL = "https://api.appnexus.com/batch-segment"
			log.Println("apiUPL: ", apiURL)
		}
	*/

	apiURL += "?member_id=" + strconv.Itoa(int(memberID))

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
		log.Println("buf:", string(buf))
		return nil, err
	}

	if errResponse.Response.Error != "" {
		return nil, fmt.Errorf("%s:%s", errResponse.Response.ErrorId, errResponse.Response.Error)
	}
	var bsResponse BatchSegmentResponse

	if err := json.Unmarshal(buf, &bsResponse); err != nil {
		log.Println("buf:", string(buf))
		return nil, err
	}

	log.Println("jobs:", len(bsResponse.Response.BatchSegmentUploadJob))
	log.Println("COMPLETE_time: ", bsResponse.Response.BatchSegmentUploadJob[0].CompletedTime)

	return bsResponse.Response.BatchSegmentUploadJob, nil
}
