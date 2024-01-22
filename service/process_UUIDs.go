package service

import (
	"errors"
	"log"
	"strconv"
	"strings"
)

func validateUUID(uuid string) ([]string, error) {
	var err error
	var errList []string

	log.Println("------------VALIDATE UUID-------------")
	sections, err := parseUUID(uuid)
	if err != nil {
		log.Println("Parsing Err: ", err)
		errList = append(errList, err.Error())
		log.Println("sections len: ", len(sections))
		for i := 0; i < len(errList); i++ {
			log.Println(errList[i])
			return errList, err

		}
		return errList, err
	}

	s := "123!4567"
	n, err := strconv.ParseUint(s, 16, 64)
	if err != nil {
		log.Println("s is not valid", err)
	}
	log.Println("n = ", n)
	log.Println("len s = ", len(s))
	return errList, err
}

func parseUUID(uu string) ([]string, error) {
	var sections []string
	var uuid string
	var found bool
	var hyphen int
	var err error

	log.Println("----------PARSE UUID----------------")

	uuid = uu
	log.Println("UUID STRING: ", uuid)
	log.Println("len of hyphens: ", strings.Count(uuid, "-"))
	hyphenQuantity := strings.Count(uuid, "-")
	if hyphenQuantity == 4 {
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
			if err != nil {
				return sections, err
			}
		}
		log.Println("UUID AFTER FOR: ", uuid)
		sections = append(sections, uuid)

		for i := 0; i < len(sections); i++ {
			log.Println(i, ". ", sections[i])
		}
		log.Println("Sections len: ", len(sections))
	} else {
		errString := "There is only " + strconv.Itoa(hyphenQuantity) + " hyphens. Must be 4!"
		errQnt := errors.New(errString)
		err = errQnt
		log.Println(errQnt)

	}

	return sections, err
}
