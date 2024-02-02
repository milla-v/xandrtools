all:
	go build xandrtools/cmd/xandrtools

run: xandrtools
	go build xandrtools/cmd/xandrtools
	DEBUG_ADDR=127.0.0.1:9001 ./xandrtools
