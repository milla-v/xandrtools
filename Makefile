all:
	go build -o /tmp/xandrtools xandrtools/cmd/xandrtools

run:
	go build -o /tmp/xandrtools xandrtools/cmd/xandrtools
	DEBUG_ADDR=127.0.0.1:9970 /tmp/xandrtools

cert:
	go run /usr/local/go/src/crypto/tls/generate_cert.go -host $HOST
	mkdir -p ~/.config/xandrtools
	mv cert.pem key.pem ~/.config/xandrtools
