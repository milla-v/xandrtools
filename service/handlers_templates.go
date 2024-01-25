package service

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var t *template.Template

func handleStyle(w http.ResponseWriter, r *http.Request) {
	log.Println("handleStyle")
	fs := http.FS(cssFiles)
	http.FileServer(fs).ServeHTTP(w, r)
}

func handlePng(w http.ResponseWriter, r *http.Request) {
	log.Println("handlePng", r.URL.Path)
	fs := http.FS(pngFiles)
	http.FileServer(fs).ServeHTTP(w, r)
}

func handleXandrtools(w http.ResponseWriter, r *http.Request) {
	type data struct {
		XUID       int64
		Validation xandr
		Errs       bool
		ValUUID    uuid
		SecOne     string
		SecTwo     string
		SecThree   string
		SecFour    string
		SecFive    string
	}
	var d data
	var err error
	d.Errs = false
	var s string

	if r.Method == "POST" {
		switch r.RequestURI {
		case "/validate?type=xandrid&id=":
			s = r.FormValue("xuid")
			if s == "" {
				fmt.Println("empty link")
			}
			d.Validation = processXandrUID(s)
			fmt.Println("---------------------------")
			log.Println("len errs: ", len(d.Validation.ErrList))
			if len(d.Validation.ErrList) > 0 {
				d.Errs = true
			}
			log.Println("errs: ", d.Errs)

		case "/validate?type=uuid&id=":
			s = r.FormValue("uuid")
			if s == "" {
				fmt.Println("empty link")
			}
			d.ValUUID, err = validateUUID(s)
			if err != nil {
				log.Println("ValUUD err: ", len(d.ValUUID.ErrMsg))
				log.Println("ErrSecNum = ", d.ValUUID.ErrSecNum)
			}
			if len(d.ValUUID.Sections) > 0 {
				d.SecOne = d.ValUUID.Sections[0]
				d.SecTwo = d.ValUUID.Sections[1]
				d.SecThree = d.ValUUID.Sections[2]
				d.SecFour = d.ValUUID.Sections[3]
				d.SecFive = d.ValUUID.Sections[4]
			}
			log.Println("errmsg: ", d.ValUUID.ErrMsg)
		}
	}

	if r.Method == "GET" && r.URL.Path != "/" {
		if r.URL.Path == "/validate" && r.URL.Query().Get("type") == "xandrid" {
			log.Println("VALIDATE TYPE: ", r.URL.Query().Get("type"))
			id := r.URL.Query().Get("id")
			log.Println("XandrID = ", id)
			d.Validation = processXandrUID(id)

			log.Println("len errs: ", len(d.Validation.ErrList))
			if len(d.Validation.ErrList) > 0 {
				d.Errs = true
			}
			log.Println("errs: ", d.Errs)
		}
		if r.URL.Path == "/validate" && r.URL.Query().Get("type") == "uuid" {
			log.Println("VALIDATE TYPE: ", r.URL.Query().Get("type"))
			id := r.URL.Query().Get("id")
			log.Println("UUID = ", id)
			d.ValUUID, err = validateUUID(s)
			if err != nil {
				log.Println("ValUUD err: ", len(d.ValUUID.ErrMsg))
				log.Println("ErrSecNum = ", d.ValUUID.ErrSecNum)
			}
			if len(d.ValUUID.Sections) > 0 {
				d.SecOne = d.ValUUID.Sections[0]
				d.SecTwo = d.ValUUID.Sections[1]
				d.SecThree = d.ValUUID.Sections[2]
				d.SecFour = d.ValUUID.Sections[3]
				d.SecFive = d.ValUUID.Sections[4]
			}
			log.Println("errmsg: ", d.ValUUID.ErrMsg)
		}
	}

	/*
		switch r.Method {
		case "POST":

			s := r.FormValue("xuid")
			log.Println("r.RequestURL: ", r.RequestURI)
			log.Println("R: ", r.Form)

			if s == "" {
				fmt.Println("empty link")
			}
			d.Validation = processXandrUID(s)
			fmt.Println("----------------------------")

		}
	*/
	//test := "123e4567-e89b-12d3-a456-426614174000"
	/*
		d.ValUUID.ErrMsg, d.ValUUID.ErrSecNum, err = validateUUID(test)
		if err != nil {
			log.Println("ValUUD err: ", len(d.ValUUID.ErrMsg))
			log.Println("ErrSecNum = ", d.ValUUID.ErrSecNum)
		}
	*/

	if err := t.ExecuteTemplate(w, "xandrtools.html", d); err != nil {
		log.Println(err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
}
