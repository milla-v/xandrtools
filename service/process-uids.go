package service

import (
	"log"
	"strconv"
	"strings"
)

func processXandrUID(str string) xandr {
	var err error
	var xuid int64
	var xr xandr
	xr.ErrList = nil

	negative := "Negative number"
	letters := "ID must contain numbers only"
	zero := "Cannot start from 0"
	large := "UserID value is out of range"
	empty := "Empty the Xanrd userID field"

	log.Println("str[0]: ", string(str[0]))

	xuid, err = strconv.ParseInt(str, 10, 64)
	if err != nil {
		if strings.Contains(err.Error(), "out of range") {
			xr.ErrList = append(xr.ErrList, large)
			log.Println(large)
		}
		if strings.Contains(err.Error(), "invalid syntax") {
			xr.ErrList = append(xr.ErrList, letters)
			log.Println(letters)
		} else {
			log.Println("other err: ", err)
		}
	} else if string(str[0]) == "0" {
		xr.ErrList = append(xr.ErrList, zero)
		log.Println(zero)
	} else if xuid == 0 {
		xr.ErrList = append(xr.ErrList, empty)
		log.Println(empty)
	} else {
		log.Println("xuid:", xuid)
		if xuid <= 0 {
			xr.ErrList = append(xr.ErrList, negative)
			log.Println("negative")
		}

	}
	for i := 0; i < len(xr.ErrList); i++ {
		log.Println(i, ". ", xr.ErrList[i])
	}
	return xr
}
