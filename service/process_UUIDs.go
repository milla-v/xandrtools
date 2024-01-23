package service

import (
	"errors"
	"log"
	"strconv"
	"strings"
)

func validateUUID(uuid string) (string, error) {
	var err error
	var errMsg string

	log.Println("------------VALIDATE UUID-------------")
	sections, err := parseUUID(uuid)
	if err != nil {
		log.Println("Parsing Err: ", err)
		errMsg = err.Error()
		return errMsg, err
	}
	//check if hexadecimal
	if len(sections) > 0 {
		for i := 0; i < len(sections); i++ {
			_, err := strconv.ParseInt(sections[i], 16, 64)
			if err != nil {
				errMsg = notHex
				log.Println("errMsg ", errMsg)
				return errMsg, err
			}
		}
	}

	return errMsg, err
}

func parseUUID(uu string) ([]string, error) {
	var sections []string
	var uuid string
	var err error

	log.Println("----------PARSE UUID----------------")

	uuid = uu
	log.Println("UUID STRING: ", uuid)
	log.Println("len of hyphens: ", strings.Count(uuid, "-"))
	hyphenQuantity := strings.Count(uuid, "-")
	if hyphenQuantity == 4 {
		sections = strings.Split(uuid, "-")
		log.Println("SPLIT res len: ", len(sections))
		for i := 0; i < len(sections); i++ {
			log.Println(i, ". ", sections[i])
		}
		if err != nil {
			log.Println("slpit err: ", err)
			return sections, err
		}
	} else {
		errString := "There is only " + strconv.Itoa(hyphenQuantity) + " hyphens. Must be 4!"
		errQnt := errors.New(errString)
		err = errQnt
		log.Println(errQnt)

	}
	log.Println("I am here")
	return sections, err
}
