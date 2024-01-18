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
		XUID        int64
		XuidErrList []string
	}
	var d data

	switch r.Method {
	case "POST":
		s := r.FormValue("xuid")
		if s == "" {
			fmt.Println("empty link")
		}
		d.XuidErrList = processXandrUID(s)
		fmt.Println("----------------------------")

	}

	if err := t.ExecuteTemplate(w, "xandrtools.html", d); err != nil {
		log.Println(err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
}
