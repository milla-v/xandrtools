package client

import (
	"crypto/tls"
	"errors"
	"log"
	"net/http"
	"strings"
)

// get  backend and debug adress string
// returns api url to login and get bss jobs
// returns err if apiURL is empty
func getApiURL(backend string, debugAddress string) (string, error) {
	var apiURL string
	var err error
	addr := debugAddress
	var isDebug bool
	var port string
	log.Println("Backend: ", backend)
	log.Println("debug address: ", debugAddress)
	if debugAddress != "" {
		isDebug = true
		port = strings.TrimPrefix(addr, "127.0.0.1")
	}
	log.Println("port: ", port)

	switch {
	case backend == "simulator" && isDebug == true:
		apiURL = "http://localhost" + port + "/xandrsim/"
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
