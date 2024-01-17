package main

// chat command runs a chat server.
package main

import (
	"flag"
	"fmt"
	"log"

	"xandrtools/service"

	_ "time/tzdata"
)

var version = flag.Bool("version", false, "Print version")

// Version is set by linker
var Version string

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()
	if *version {
		fmt.Println("version:", Version)
		return
	}

	service.Run()
}