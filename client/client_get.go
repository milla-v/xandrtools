package client

import (
	"crypto/tls"
	"errors"
	"log"
	"net/http"
	"os"
)

// get  backend string
// returns api url to login and get bss jobs
// returns err if apiURL is empty
func getApiURL(backend string) (string, error) {
	var apiURL string
	var err error
	debugAddr := os.Getenv("DEBUG_ADDR")
	var isDebug bool
	if debugAddr != "" {
		isDebug = true
	}

	switch {
	case backend == "simulator" && isDebug == true:
		apiURL = "http://" + debugAddr + "/xandrsim/"
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		log.Println("apiUPL: ", apiURL)
	case backend == "simulator" && isDebug == false:
		apiURL = "https://xandrtools.com/xandrsim/"
		log.Println("apiUPL: ", apiURL)
	case backend == "xandr" && isDebug == false:
		apiURL = "https://api.appnexus.com/"
		log.Println("apiUPL: ", apiURL)
	}
	if apiURL == "" {
		err = errors.New("Api url is empty")
	}
	log.Println("Error: ", err)
	return apiURL, err
}
