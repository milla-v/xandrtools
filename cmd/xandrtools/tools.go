package main

import (
	"log"
	"runtime/debug"
)

func getTag() (string, error) {
	var err error
	info, _ := debug.ReadBuildInfo()
	log.Println("info: ", info)
	deps, _ := debug.ReadBuildInfo()
	//log.Println("dep 0: ", deps.Deps[0].Path, deps.Deps[0].Version)
	var tag string
	for i := 0; i < len(deps.Deps); i++ {
		//log.Println(i, " : ", "Path: ", deps.Deps[i].Path, " | Version: ", deps.Deps[i].Version)
		if deps.Deps[i].Path == "github.com/milla-v/xandr" {
			tag = deps.Deps[i].Version
			//log.Println("FOUND TAG: ", tag)
		}
	}

	return tag, err
}
