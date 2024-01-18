package service

import (
	"log"
	"strconv"
)

func processXandrUID(str string) []string {

	var xuid int64
	var xuidErrList []string
	xuidErrList = nil

	negative := "Negative number"
	//zero := "cannot start from 0"

	log.Println("str[0]: ", string(str[0]))

	xuid, _ = strconv.ParseInt(str, 10, 64)
	if xuid == 0 {
		return xuidErrList
		log.Println("empty the Xanrd userID field")
	}
	log.Println("xuid:", xuid)
	if xuid <= 0 {
		xuidErrList = append(xuidErrList, negative)
	}

	for i := 0; i < len(xuidErrList); i++ {
		log.Println(i, ". ", xuidErrList[i])
	}
	return xuidErrList
}
