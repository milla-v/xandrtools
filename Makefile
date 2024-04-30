all:
	go build -o /tmp/xandrtools xandrtools/cmd/xandrtools

run:
	go build -o /tmp/xandrtools xandrtools/cmd/xandrtools
	DEBUG_ADDR=127.0.0.1:9001 /tmp/xandrtools
