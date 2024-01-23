package service

import (
	"errors"
	"log"
	"strconv"
	"strings"
)

func validateUUID(str string) (uuid, error) {
	var err error
	var u uuid

	log.Println("------------VALIDATE UUID-------------")
	u.UUID = str
	u.Sections, err = parseUUID(str)
	if err != nil {
		log.Println("Parsing Err: ", err)
		u.ErrMsg = err.Error()
		u.SectionsExist = false
		log.Println("Sections Exist: ", u.SectionsExist)
		return u, err
	}
	//check if hexadecimal
	log.Println("Len u.Sections = ", len(u.Sections))
	if len(u.Sections) > 0 {
		u.SectionsExist = true
		for i := 0; i < len(u.Sections); i++ {
			_, err := strconv.ParseInt(u.Sections[i], 16, 64)
			if err != nil {
				//add 1 to don't get 0 section if error
				u.ErrSecNum = i + 1
				u.ErrMsg = "Section " + strconv.Itoa(u.ErrSecNum) + " " + notHex
				log.Println("errMsg ", u.ErrMsg)
				log.Println("eerSecNum: ", u.ErrSecNum)
				return u, err
			}
		}
	}
	log.Println("Sections Exist: ", u.SectionsExist)
	return u, err
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
		errString := "There is only " + strconv.Itoa(hyphenQuantity) + " hyphens. Must be 4! Ex.: XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
		errQnt := errors.New(errString)
		err = errQnt
		log.Println(errQnt)

	}
	log.Println("I am here")
	return sections, err
}
