// Package service implements chat http service.
package service

//go:generate go run ../cmd/chatembed/chatembed.go

var proxyXandrtools = func() *httputil.ReverseProxy {
	u, err := url.Parse("http://zero:80/xandrtools/")
	if err != nil {
		log.Fatal(err)
	}
	return httputil.NewSingleHostReverseProxy(u)
}()

// Run starts a chat http server on address (host:port)
func Run() {
	var err error

	cfg.Readenv()
	domcfg, err = loadDomainsConfig()
	if err != nil {
		log.Println("cannot load domains config:", err)
	} else {
		log.Printf("domains loaded: %+v", domcfg)
	}

	log.Printf("chat version: %s, date: %s\n", version, date)
	log.Println("starting server on https://" + cfg.Domain + ":" + cfg.Port + "/")

	loc = mustLoadLocation()

	connectChan = make(chan *client)
	connectedChan = make(chan *client, 100)
	disconnectChan = make(chan *client, 100)
	broadcastChan = make(chan *message, 100)
	historyFile, err = os.OpenFile(cfg.WorkDir+"history.html", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mux.Handle("xandrtools.com/", proxyXandrtools)

	_, err = os.Stat("cert.pem")
	if err == nil {
		log.Println("cert.pem exists. Starting dev server")
		hosts := addDynamicHosts(nil)
		log.Printf("dynamic hosts: %s", hosts)
		startDevServer(mux)
		return
	}

	log.Println("cert.pem does not exist. Starting prod server with autocert")
}

func startDevServer(h http.Handler) {
	var s *http.Server

	if strings.HasPrefix(cfg.Domain, "127.0.0.1") {
		s = &http.Server{
			Addr:    cfg.Domain + ":" + cfg.Port,
			Handler: h,
		}
	} else {
		s = &http.Server{
			Addr:    ":" + cfg.Port,
			Handler: h,
		}
	}
	log.Fatal(s.ListenAndServeTLS("cert.pem", "key.pem"))
}
