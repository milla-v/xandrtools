package service

import (
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
		XUID             int64
		ValidationResult xandr
		Errs             bool
		ValUUID          uuid
		SecOne           string //section one of uuid
		SecTwo           string //section two of uuid
		SecThree         string //section three of uuid
		SecFour          string //section four of uuid
		SecFive          string //section five of uuid
	}
	var d data
	var err error
	d.Errs = false

	//1. input name = "type" value = "xandrid"
	if r.URL.Query().Get("type") == "xandrid" {
		id := r.URL.Query().Get("id")
		d.ValidationResult = processXandrUID(id)
		if len(d.ValidationResult.ErrList) > 0 {
			d.Errs = true
		}
	}
	//2. input name = "type" value = "uuid"
	if r.URL.Query().Get("type") == "uuid" {
		id := r.URL.Query().Get("id")
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
		ID            string
		ShowText      bool
		Link          string
		InitScript    template.JS
		GeneratedText string
		Seps          separators
		SegError      string //holds segment errors
		SepError      string //separator errors

	}

	var err error
	var d data
	d.ShowText = false

	d.Seps.Sep1 = r.URL.Query().Get("sep_1")
	d.Seps.Sep2 = r.URL.Query().Get("sep_2")
	d.Seps.Sep3 = r.URL.Query().Get("sep_3")
	d.Seps.Sep4 = r.URL.Query().Get("sep_4")
	d.Seps.Sep5 = r.URL.Query().Get("sep_5")

	//set default separator
	setDefaultSeparators(&d.Seps)

	//check separators
	if err := checkSeparators(d.Seps); err != nil {
		d.SepError = err.Error()
	}

	//check sf value
	sf := r.URL.Query().Get("sf")
	segmentFields := strings.Split(sf, "-")

	// checks segments
	d.SegError, err = checkSegments(segmentFields)
	if err != nil {
		d.SegError = err.Error()
		log.Println("d.SegError error: ", d.SegError)
	}

	var js string
	for _, f := range segmentFields {
		id := "'" + strings.ToLower(f) + "'"
		js += "var checkBox = document.getElementById(" + id + ");\n"
		js += "checkBox.checked = true;\n"
		js += "checkField(" + id + ");\n"
	}
	d.InitScript = template.JS(js)

	//generate text sample
	if len(d.SegError) == 0 && sf != "" && len(d.SepError) == 0 {
		d.ShowText = true
		d.GeneratedText = generateSample(segmentFields, d.Seps)
	}
	if err := t.ExecuteTemplate(w, "textGenerator.html", d); err != nil {
		log.Println(err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
}
