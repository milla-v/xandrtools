package service

import (
	"log"
	"strconv"
	"strings"
)

func validateUUID() {
	s := "123!4567"
	n, err := strconv.ParseUint(s, 16, 64)
	if err != nil {
		log.Println("s is not valid", err)
	}
	log.Println("n = ", n)
	log.Println("len s = ", len(s))
}

func parseUUID(uu string) {
	var sections []string
	var uuid string
	var found bool
	var hyphen int
	uuid = uu
	log.Println("UUID STRING: ", uuid)
	for i := 0; i < 4; i++ {
		//find, copy and delete with hyphen first 4 sections of uuid
		log.Println("-------------------", i, "--------------------")
		hyphen = strings.Index(uuid, "-")
		log.Println("First hyphen position: ", hyphen)
		sec := uuid[:hyphen]
		log.Println("sec1: ", sec)
		sections = append(sections, sec)
		uuid, found = strings.CutPrefix(uuid, sec+"-")
		log.Println("uuid after delete first section: ", uuid, "found: ", found)
	}

	/*
		//find, copy and delete with hyphen first section of id
		hyphen := strings.Index(uuid, "-")
		log.Println("First hyphen position: ", hyphen)
		sec := uuid[:hyphen]
		log.Println("sec1: ", sec)
		sections = append(sections, sec)
		uuid, found = strings.CutPrefix(uuid, sec+"-")
		log.Println("uuid after delete first section: ", uuid, "found: ", found)

		//find, copy and delete with hyphen first section of id
		hyphen = strings.Index(uuid, "-")
		log.Println("Second hyphen p[osition: ", hyphen)
		sec = uuid[:hyphen]
		log.Println("sec: ", sec)
		sections = append(sections, sec)
		uuid, found = strings.CutPrefix(uuid, sec+"-")
		log.Println("uuid after delete section: ", uuid, "found: ", found)
	*/
	log.Println("UUID AFTER FOR: ", uuid)
	sections = append(sections, uuid)

	for i := 0; i < len(sections); i++ {
		log.Println(i, ". ", sections[i])
	}
	return
}
