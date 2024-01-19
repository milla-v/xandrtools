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
	xr.UserID = 0

	negative := "ID cannot be negative"
	letters := "ID must contain numbers only"
	zero := "Cannot start from 0"
	large := "UserID value is out of range"
	empty := "Empty the Xanrd userID field"
	valid := "ID validation completed. No errors or warnings founded."

	log.Println("str[0]: ", string(str[0]))

	xuid, err = strconv.ParseInt(str, 10, 64)
	if err != nil {
		if strings.Contains(err.Error(), "out of range") {
			xr.ErrList = append(xr.ErrList, large)
			log.Println(large)
			xr.WrongUserID = str
		}
		if strings.Contains(err.Error(), "invalid syntax") {
			xr.ErrList = append(xr.ErrList, letters)
			log.Println(letters)
			xr.WrongUserID = str
		} else {
			log.Println("other err: ", err)
		}
	} else if string(str[0]) == "0" {
		xr.ErrList = append(xr.ErrList, zero)
		log.Println(zero)
		xr.WrongUserID = str
	} else if xuid == 0 {
		xr.ErrList = append(xr.ErrList, empty)
		log.Println(empty)
		xr.WrongUserID = "Empty input"
	} else if xuid <= 0 {
		xr.ErrList = append(xr.ErrList, negative)
		log.Println("negative")
		xr.WrongUserID = str
	} else {
		xr.ValidMsg = valid
		xr.UserID = xuid
	}
	for i := 0; i < len(xr.ErrList); i++ {
		log.Println(i, ". ", xr.ErrList[i])
	}
	return xr
}
