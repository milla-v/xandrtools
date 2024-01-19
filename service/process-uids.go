package service

import (
	"log"
	"strconv"
)

func processXandrUID(str string) []string {
	var err error
	var xuid int64
	var xuidErrList []string
	xuidErrList = nil

	negative := "Negative number"
	letters := "ID must contain numbers only"
	zero := "cannot start from 0"

	log.Println("str[0]: ", string(str[0]))

	xuid, err = strconv.ParseInt(str, 10, 64)
	if err != nil {
		log.Println("xuid parcing err: ", err)
		xuidErrList = append(xuidErrList, letters)
	} else if string(str[0]) == "0" {
		log.Println(zero)
	} else if xuid == 0 {
		log.Println("empty the Xanrd userID field")
	} else {
		log.Println("xuid:", xuid)
		if xuid <= 0 {
			xuidErrList = append(xuidErrList, negative)
		}
	}
	for i := 0; i < len(xuidErrList); i++ {
		log.Println(i, ". ", xuidErrList[i])
	}
	return xuidErrList
}
