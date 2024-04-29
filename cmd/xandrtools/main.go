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

func main() {
	var err error
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()
	if *version {
		fmt.Println("version:", version)
		return
	}

	service.Version, err = getTag()
	if err != nil {
		log.Println("getTag err: ", err)
		return
	}
	service.Run()

}
