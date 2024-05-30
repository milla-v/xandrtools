package service

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/milla-v/xandr/bss/xgen"

	"xandrtools/simulator"
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
	log.Println("start xandrtools")
	type data struct {
		XandrVersion string
		VCS          Vcs
	}
	var d data

	d.XandrVersion = Version
	d.VCS.RevisionFull = VcsInfo.RevisionFull
	d.VCS.RevisionShort = VcsInfo.RevisionShort
	d.VCS.Modified = VcsInfo.Modified

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
		InitScript    template.JS
		GeneratedText string
		Seps          separators
		GenError      string //error from xgen library
		GenText       string
		XandrVersion  string
		VCS           Vcs
	}

	var err error
	var d data
	d.ShowText = false

	d.XandrVersion = Version
	d.VCS.RevisionFull = VcsInfo.RevisionFull
	d.VCS.RevisionShort = VcsInfo.RevisionShort
	d.VCS.Modified = VcsInfo.Modified

	sfs := r.URL.Query().Get("sf")
	log.Println("sfs: ", sfs)
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

	log.Println("params => segmentFields => ", params.SegmentFields)

	if len(d.GenError) == 0 && sfs != "" {
		d.ShowText = true
		d.GeneratedText, err = generateSample(&params)
		if err != nil {
			d.GenError = err.Error()
		}
	}

	d.Seps.Sep1 = r.URL.Query().Get("sep_1")
	d.Seps.Sep2 = r.URL.Query().Get("sep_2")
	d.Seps.Sep3 = r.URL.Query().Get("sep_3")
	d.Seps.Sep4 = r.URL.Query().Get("sep_4")
	d.Seps.Sep5 = r.URL.Query().Get("sep_5")

	setDefaultSeparators(&d.Seps)

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
		XandrVersion     string
		VCS              Vcs
	}
	var d data
	var err error
	d.Errs = false

	d.XandrVersion = Version
	d.VCS.RevisionFull = VcsInfo.RevisionFull
	d.VCS.RevisionShort = VcsInfo.RevisionShort
	d.VCS.Modified = VcsInfo.Modified

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

func handleBssTroubleShooter(w http.ResponseWriter, r *http.Request) {
	log.Println("start bss trouble shooter")
	type data struct {
		XandrVersion string
		VCS          Vcs

		Auth simulator.AuthRequest
		User simulator.UserData

		Communicator string
		Token        string
		Backend      string
	}
	var d data

	d.XandrVersion = Version
	d.VCS.RevisionFull = VcsInfo.RevisionFull
	d.VCS.RevisionShort = VcsInfo.RevisionShort
	d.VCS.Modified = VcsInfo.Modified

	log.Println("Method before submit", r.Method)
	//check if username and password not empty
	//get username and passssword
	//d.UserName = r.FormValue("username")

	//authentication request
	d.Auth.Auth.Username = r.FormValue("username")
	d.Auth.Auth.Password = r.FormValue("password")
	d.Backend = r.FormValue("backend")
	if d.Auth.Auth.Username != "" && d.Auth.Auth.Password != "" {
		//sendRequest JSON with password and username
		buf, err := json.MarshalIndent(d.Auth, "\t", "\t")
		if err != nil {
			log.Println("Marshal err: ", err)
			return
		}

		log.Println("json:", string(buf))

		var apiURL = "https://api.appnexus.com/auth"
		if d.Backend == "simulator" {
			apiURL = "http://127.0.0.1:9970/xandrsim/auth"
		}

		resp, err := http.Post(apiURL, "application/json", bytes.NewReader(buf))
		if err != nil {
			log.Println("Post responce err: ", err)
			return
		}

		buff, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("read responce body err: ", err)
			return
		}
		defer resp.Body.Close()

		var u simulator.AuthResponse
		log.Printf("status:%s body: %s", resp.Status, string(buff))
		if err := json.Unmarshal(buf, &u); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Println(u.Response.Status)
		//fill in an user datas
		//d.User.Username = d.Auth.Auth.Username

	}
	log.Println("d.User.Username: ", d.Auth.Auth.Username, " | ", "d.User.Password: ", d.Auth.Auth.Password, " | ", "d.Token: ", d.Token)

	//get token, check if token not empty
	//d.Token = simulator.

	if err := t.ExecuteTemplate(w, "bsstroubleshooter.html", d); err != nil {
		log.Println("Execute bsstroubleshooter.html ", err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
}
