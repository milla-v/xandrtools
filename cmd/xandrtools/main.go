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
	err := getTag()
	if err != nil {
		log.Println("getTag err: ", err)
		return
	}
	log.Println("TAG Vesion: ", Version)

	service.Run()

}
