// Package service implements chat http service.
package service

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"golang.org/x/crypto/acme/autocert"

	"xandrtools/simulator"
)

//go:embed templates/*.html
var content embed.FS

//go:embed templates/*.css
var cssFiles embed.FS

//go:embed templates/*.png
var pngFiles embed.FS

var Version string
var VcsInfo Vcs

var proxyXandrtools = func() *httputil.ReverseProxy {
	u, err := url.Parse("http://zero:80/xandrtools/")
	if err != nil {
		log.Fatal(err)
	}
	return httputil.NewSingleHostReverseProxy(u)
}()

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "home page")
	log.Printf(r.URL.Path)
}

// Run starts http or https server
func Run() {
	t = template.Must(template.ParseFS(content,
		"templates/xandrtools.html",
		"templates/bsstroubleshooter.html",
		"templates/textGenerator.html",
		"templates/validators.html"))

	mux := http.NewServeMux()

	//	mux.Handle("xandrtools.com/", proxyXandrtools)
	mux.HandleFunc("/", handleXandrtools)
	mux.HandleFunc("/bsstroubleshooter", handleBssTroubleShooter)
	mux.HandleFunc("/textGenerator", handleTextGenerator)
	mux.HandleFunc("/validators", handleValidators)

	mux.HandleFunc("/xandrsim/auth", simulator.HandleAuthentication)
	mux.HandleFunc("/xandrsim/batch-segment", simulator.HandleBatchSegment)

	mux.HandleFunc("/templates/styles.css", handleStyle)
	mux.HandleFunc("/templates/copy_btn.png", handlePng)

	addr := os.Getenv("DEBUG_ADDR")
	if addr != "" {
		startDevServer(mux, addr)
	} else {
		startProdServer(mux)
	}
}

// startProdServer get the certificate for domain and starts a https server.
func startProdServer(h http.Handler) {
	hosts := []string{
		"xandrtools.com",
	}
	log.Println("starting prod server")
	l := autocert.NewListener(hosts...)
	err := http.Serve(l, h)
	log.Fatal(err)
}

// startDevServer starts http server on localhost.
func startDevServer(h http.Handler, addr string) {
	log.Println("starting dev server on", addr)
	s := &http.Server{
		Addr:    addr,
		Handler: h,
	}

	certPath := os.Getenv("HOME") + "/.config/xandrtools/cert.pem"

	_, err := os.Stat(certPath)
	if err == nil {
		log.Println("local certificate found")
		keyPath := os.Getenv("HOME") + "/.config/xandrtools/key.pem"
		err = s.ListenAndServeTLS(certPath, keyPath)
	} else {
		err = s.ListenAndServe()
	}

	if err != nil {
		fmt.Println("Error:", err)
	}
}
