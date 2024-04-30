package main

import (
	"errors"
	"log"
	"runtime/debug"

	"xandrtools/service"
)

func getTag() (string, error) {
	var err error
	info, _ := debug.ReadBuildInfo()
	log.Println("info: ", info)
	deps, _ := debug.ReadBuildInfo()
	//log.Println("build 0: ", deps.Settings[0].Key, deps.Settings[0].Value)
	var tag string
	for i := 0; i < len(deps.Deps); i++ {
		//log.Println(i, " : ", "Path: ", deps.Deps[i].Path, " | Version: ", deps.Deps[i].Version)
		if deps.Deps[i].Path == "github.com/milla-v/xandr" {
			tag = deps.Deps[i].Version
		}
	}
	if len(tag) == 0 {
		err = errors.New("Tag is NULL")
	}
	return tag, err
}

func getVCS() (service.Vcs, error) {
	var err error
	var vcsInfo service.Vcs
	var modified string
	vcsInfo.Modified = false

	deps, _ := debug.ReadBuildInfo()
	for i := 0; i < len(deps.Settings); i++ {
		log.Println("SETTINGS: ", deps.Settings[i].Key, " : ", deps.Settings[i].Value)
		if deps.Settings[i].Key == "vcs.revision" {
			log.Println("vcs.revision found!")
			vcsInfo.RevisionFull = deps.Settings[i].Value
			vcsInfo.RevisionShort = deps.Settings[i].Value[:7]
		}
		if deps.Settings[i].Key == "vcs.modified" {
			modified = deps.Settings[i].Value
		}
	}
	if len(vcsInfo.RevisionFull) == 0 {
		err = errors.New("vcs.revision full is empty")
	}
	if len(vcsInfo.RevisionShort) == 0 {
		err = errors.New("vcs.revision short is empty")
	}
	if len(modified) == 0 {
		err = errors.New("vcs.modified is empty")
	}
	if modified == "true" {
		vcsInfo.Modified = true
	}

	log.Println("REVISION FULL: ", vcsInfo.RevisionFull)
	log.Println("REVISION SHORT: ", vcsInfo.RevisionShort)
	log.Println("MODIFIED: ", vcsInfo.Modified)
	return vcsInfo, err
}
