package main

import (
	"log"
	"runtime/debug"
)

type vcs struct {
	RevisionFull  string
	RevisionShort string
	Modified      bool
}

func getBuildInfo() (string, error) {
	var err error
	var vcsInfo vcs
	vcsInfo.Modified = false
	info, _ := debug.ReadBuildInfo()
	log.Println("info: ", info)
	deps, _ := debug.ReadBuildInfo()
	log.Println("build 0: ", deps.Settings[0].Key, deps.Settings[0].Value)
	var tag string
	for i := 0; i < len(deps.Deps); i++ {
		//log.Println(i, " : ", "Path: ", deps.Deps[i].Path, " | Version: ", deps.Deps[i].Version)
		if deps.Deps[i].Path == "github.com/milla-v/xandr" {
			tag = deps.Deps[i].Version
		}
	}
	for i := 0; i < len(deps.Settings); i++ {
		log.Println("SETTINGS: ", deps.Settings[i].Key, " : ", deps.Settings[i].Value)
		if deps.Settings[i].Key == "vcs.revision" {
			log.Println("vcs.revision found!")
			vcsInfo.RevisionFull = deps.Settings[i].Value
			vcsInfo.RevisionShort = deps.Settings[i].Value[:7]
		}
		if deps.Settings[i].Key == "vcs.modified" && deps.Settings[i].Value == "true" {
			vcsInfo.Modified = true
		}
	}

	log.Println("REVISION FULL: ", vcsInfo.RevisionFull)
	log.Println("REVISION SHORT: ", vcsInfo.RevisionShort)
	log.Println("MODIFIED: ", vcsInfo.Modified)

	return tag, err
}
