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

func handleXandrtools(w http.ResponseWriter, r *http.Request) {
	type data struct {
		XUID       int64
		Validation xandr
		Errs       bool
		ValUUID    uuid
	}
	var d data
	var err error
	d.Errs = false
	switch r.Method {
	case "POST":
		s := r.FormValue("xuid")
		if s == "" {
			fmt.Println("empty link")
		}
		d.Validation = processXandrUID(s)
		fmt.Println("----------------------------")

	}
	test := "123e4567-e89b-12d3-a456-426614174000"
	d.ValUUID.ErrMsg, err = validateUUID(test)
	if err != nil {
		log.Println("ValUUD err: ", len(d.ValUUID.ErrMsg))
	}

	log.Println("len errs: ", len(d.Validation.ErrList))
	if len(d.Validation.ErrList) > 0 {
		d.Errs = true
	}
	log.Println("errs: ", d.Errs)

	if err := t.ExecuteTemplate(w, "xandrtools.html", d); err != nil {
		log.Println(err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
}
