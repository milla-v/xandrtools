package service

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/milla-v/xandr/bss/xgen"
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
	if err := t.ExecuteTemplate(w, "xandrtools.html", nil); err != nil {
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
		InitScript    template.JS
		GeneratedText string
		Seps          separators
		GenError      string //error from xgen library
		GenText       string
	}

	var err error
	var d data
	d.ShowText = false

	sfs := r.URL.Query().Get("sf")
	str := strings.Split(sfs, "-")
	var segFields []xgen.SegmentFieldName

	for _, s := range str {
		segFields = append(segFields, xgen.SegmentFieldName(s))
	}

	//set default separator

	params := xgen.TextEncoderParameters{
		Sep1:          replaceTabs(r.URL.Query().Get("sep_1")),
		Sep2:          r.URL.Query().Get("sep_2"),
		Sep3:          r.URL.Query().Get("sep_3"),
		Sep4:          replaceTabs(r.URL.Query().Get("sep_4")),
		Sep5:          r.URL.Query().Get("sep_5"),
		SegmentFields: segFields,
	}

	//check separators, segments and return err
	_, err = xgen.NewTextEncoder(params)
	if err != nil {
		d.GenError = err.Error()
		log.Println("d.GenError = ", d.GenError)
	}

	var js string

	for _, f := range str {
		id := "'" + strings.ToLower(f) + "'"
		js += "var checkBox = document.getElementById(" + id + ");\n"
		js += "checkBox.checked = true;\n"
		js += "checkField(" + id + ");\n"
	}
	d.InitScript = template.JS(js)

	//generate text sample

	if len(d.GenError) == 0 && sfs != "" {
		d.ShowText = true
		d.GeneratedText, err = generateSample2(&params)
		if err != nil {
			d.GenError = err.Error()
		}
	}

	// old code from here

	d.Seps.Sep1 = r.URL.Query().Get("sep_1")
	d.Seps.Sep2 = r.URL.Query().Get("sep_2")
	d.Seps.Sep3 = r.URL.Query().Get("sep_3")
	d.Seps.Sep4 = r.URL.Query().Get("sep_4")
	d.Seps.Sep5 = r.URL.Query().Get("sep_5")

	setDefaultSeparators(&d.Seps)

	//generate text sample
	//	var text []string
	//	for _, s := range segFields {
	//		text = append(text, string(s))
	//	}
	//	if len(d.GenError) == 0 && sfs != "" {
	//		d.ShowText = true
	//		d.GeneratedText = generateSample(text, d.Seps)
	//	}
	/*
		if len(d.SegError) == 0 && sf != "" && len(d.SepError) == 0 {
			d.ShowText = true
			d.GeneratedText = generateSample(segmentFields, d.Seps)
		}*/
	if err := t.ExecuteTemplate(w, "textGenerator.html", d); err != nil {
		log.Println(err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
}

func handleValidators(w http.ResponseWriter, r *http.Request) {
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

	if err := t.ExecuteTemplate(w, "validators.html", d); err != nil {
		log.Println(err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
}
