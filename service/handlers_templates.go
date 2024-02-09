package service

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
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
		//1. path = validate and type = xandr
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
		//2. path = validate and type = uuid
		if r.URL.Path == "/validate" && r.URL.Query().Get("type") == "uuid" {
			log.Println("VALIDATE TYPE: ", r.URL.Query().Get("type"))
			id := r.URL.Query().Get("id")
			log.Println("UUID = ", id)
			d.ValUUID, err = validateUUID(id)
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

func handleTextGenerator(w http.ResponseWriter, r *http.Request) {
	log.Println("textGenerator page")
	type data struct {
		ID             string
		SegmentsExists bool
		Link           string
	}
	var d data
	d.SegmentsExists = false
	var segs []string

	log.Println("ID = ", d.ID)
	log.Println("initializing len of segs: ", len(segs))
	fieldsmap := make(map[int]string)

	if r.Method == "POST" {
		log.Println("r.Method = ", r.Method)
		r.ParseForm()
		//get a map of form values
		//form sorted fieldsmap
		log.Println("ID = (if POST)", d.ID)
		for k := range r.Form {
			fmt.Printf("%s value is %v\n", k, r.Form.Get(k))
			value := r.Form.Get(k)
			key, err := strconv.Atoi(strings.TrimPrefix(k, "sel_"))
			if err != nil {
				log.Println("key error :", err)
			}
			fieldsmap[key] = value
		}
		//sort by keys
		keys := make([]int, 0, len(fieldsmap))
		for k := range fieldsmap {
			keys = append(keys, k)
		}
		sort.Ints(keys)

		//get sorted array of segments
		//get id string
		for _, k := range keys {
			log.Println("Key", k, "Value", fieldsmap[k])
			//get array of segments
			segs = append(segs, fieldsmap[k])
			//form id string
			value := fieldsmap[k]
			d.ID = d.ID + value
		}
		log.Println("SEGS : ", segs)

		//to escape phohibited symbols
		d.Link = url.QueryEscape(d.ID)
		log.Println("d.Link : ", d.Link)
	}
	log.Println("ID = (after POST)", d.ID)
	if len(segs) > 0 {
		d.SegmentsExists = true
	}

	//check url and take id
	log.Println("URL: ", r.URL)
	id := r.URL.Query().Get("id")
	log.Println("GET id :", id)
	if len(id) > 0 {
		d.SegmentsExists = true
		d.ID = id
	}

	if err := t.ExecuteTemplate(w, "textGenerator.html", d); err != nil {
		log.Println(err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
}
