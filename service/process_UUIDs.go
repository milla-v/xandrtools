package service

import (
	"log"
	"strconv"
)

func validateUUID() {
	s := "123e4567"
	n, err := strconv.ParseUint(s, 16, 64)
	if err != nil {
		log.Println("s is not valid")
	}
	log.Panicln("n = ", n)
}
