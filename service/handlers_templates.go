package service

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
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
		InitScript     template.JS
		GeneratedText  string
		Seps           separators
		Errors         string
		Segs           []segments
	}
	var d data
	d.SegmentsExists = false
	var segs []string

	//log.Println("ID = ", d.ID)
	log.Println("initializing len of segs: ", len(segs))

	d.Seps.Sep1 = r.URL.Query().Get("sep_1")
	d.Seps.Sep2 = r.URL.Query().Get("sep_2")
	d.Seps.Sep3 = r.URL.Query().Get("sep_3")
	d.Seps.Sep4 = r.URL.Query().Get("sep_4")
	d.Seps.Sep5 = r.URL.Query().Get("sep_5")

	log.Println("SEPS: ", d.Seps)

	setDefaultSeparators(&d.Seps)

	//check separators
	if err := checkSeparators(d.Seps); err != nil {
		d.Errors = err.Error()
	}
	log.Println("d.Errors : ", d.Errors)
	sf := r.URL.Query().Get("sf")
	segmentFields := strings.Split(sf, "-")

	log.Println("segmentFields : ", segmentFields)

	// checks segments
	if err := checkSegments(d.Segs); err != nil {
		d.Errors = err.Error()
	}

	if sf != "" {
		d.SegmentsExists = true
	} else {
		d.Errors = "Choose at least  SEG_ID or SEG_CODE"
	}

	var js string
	for _, f := range segmentFields {
		id := "'" + strings.ToLower(f) + "'"
		js += "var checkBox = document.getElementById(" + id + ");\n"
		js += "checkBox.checked = true;\n"
		js += "checkField(" + id + ");\n"
	}
	d.InitScript = template.JS(js)

	//to escape phohibited symbols
	d.Link = r.URL.RawPath

	d.GeneratedText = generateSample(segmentFields, d.Seps)

	if err := t.ExecuteTemplate(w, "textGenerator.html", d); err != nil {
		log.Println(err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
}
