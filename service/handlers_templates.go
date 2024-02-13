package service

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	guuid "github.com/google/uuid"
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

type separators struct {
	Sep1 string
	Sep2 string
	Sep3 string
	Sep4 string
	Sep5 string
}

func handleTextGenerator(w http.ResponseWriter, r *http.Request) {
	log.Println("textGenerator page")
	type data struct {
		ID             string
		SegmentsExists bool
		Link           string
		InitScript     template.JS
		GeneratedText  string
	}
	var d data
	d.SegmentsExists = false
	var segs []string

	log.Println("ID = ", d.ID)
	log.Println("initializing len of segs: ", len(segs))

	seps := separators{
		Sep1: r.URL.Query().Get("sep_1"),
		Sep2: r.URL.Query().Get("sep_2"),
		Sep3: r.URL.Query().Get("sep_3"),
		Sep4: r.URL.Query().Get("sep_4"),
		Sep5: r.URL.Query().Get("sep_5"),
	}

	sf := r.URL.Query().Get("sf")
	segmentFields := strings.Split(sf, "-")

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

	if sf != "" {
		d.SegmentsExists = true
	}

	d.GeneratedText = generateSample(segmentFields, seps)

	if err := t.ExecuteTemplate(w, "textGenerator.html", d); err != nil {
		log.Println(err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
}

func generateSegments(segmentFields []string, seps separators, count int) (segmentsToAdd, segmentsToRem []string) {
	for i := 0; i < count; i++ {
		var segmentAdd []string
		var segmentRem []string

		for _, sf := range segmentFields {
			switch sf {
			case "SEG_ID":
				segID := 1000 + rand.Intn(1000)
				segmentAdd = append(segmentAdd, strconv.Itoa(segID))
				segmentRem = append(segmentRem, strconv.Itoa(segID+100))
			case "SEG_CODE":
				segID := 1000 + rand.Intn(1000)
				segmentAdd = append(segmentAdd, "code_"+strconv.Itoa(segID))
				segmentRem = append(segmentRem, "code_"+strconv.Itoa(segID+100))
			case "MEMBER_ID":
				segmentAdd = append(segmentAdd, "100")
				segmentRem = append(segmentRem, "100")
			case "EXPIRATION":
				segmentAdd = append(segmentAdd, "43200")
				segmentRem = append(segmentRem, "-1")
			case "VALUE":
				value := 1 + rand.Intn(5)
				segmentAdd = append(segmentAdd, strconv.Itoa(value))
				segmentRem = append(segmentRem, "0")
			}
		}

		segmentsToAdd = append(segmentsToAdd, strings.Join(segmentAdd, seps.Sep3))
		segmentsToRem = append(segmentsToRem, strings.Join(segmentRem, seps.Sep3))
	}

	return
}

type idtype struct {
	domain string
	number int
}

func generateSample(segmentFields []string, seps separators) string {
	const lineTemplate = "{UID}{SEP_1}{SEGMENTS_ADD}{SEP_4}{SEGMENTS_DEL}{SEP_5}{DOMAIN}"

	idtypes := []idtype{
		{"xandrid", 0},
		{"idfa", 3},
		{"aaid", 8},
	}

	var s string

	for _, idt := range idtypes {

		var uid string

		switch idt.domain {
		case "xandrid":
			uid = strconv.Itoa(int(rand.Int63()))
		case "idfa", "aaid":
			uid = guuid.New().String()
		default:
			log.Println("ERROR: invalid domain", idt.domain)
			continue
		}

		segmentsToAdd, segmentsToRem := generateSegments(segmentFields, seps, 2)

		var domain, sep5 string
		if idt.number != 0 {
			domain = strconv.Itoa(idt.number)
			sep5 = seps.Sep5
		}

		sr := strings.NewReplacer(
			"{UID}", uid,
			"{SEP_1}", seps.Sep1,
			"{SEGMENTS_ADD}", strings.Join(segmentsToAdd, seps.Sep2),
			"{SEP_4}", seps.Sep4,
			"{SEGMENTS_DEL}", strings.Join(segmentsToRem, seps.Sep2),
			"{SEP_5}", sep5,
			"{DOMAIN}", domain,
		)

		s += sr.Replace(lineTemplate) + "\n"
	}

	return s
}
