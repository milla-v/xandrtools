package service

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/milla-v/xandr/bss/xgen"

	"xandrtools/client"
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
		XandrVersion     string
		VCS              Vcs
		User             XandrUser
		JobList          []WebsiteBSUJ
		Backend          string
		IsJobs           bool
		IsLogin          bool
		IsLoginWithToken bool
	}
	var d data
	var err error

	d.XandrVersion = Version
	d.VCS.RevisionFull = VcsInfo.RevisionFull
	d.VCS.RevisionShort = VcsInfo.RevisionShort
	d.VCS.Modified = VcsInfo.Modified
	d.IsJobs = false
	var memberid int

	/*
		mux := http.NewServeMux()
		addr := os.Getenv("DEBUG_ADDR")
		if addr != "" {
			startDevServer(mux, addr)
		} else {
			startProdServer(mux)
		}

		log.Println("START SERVER: ", mux)
	*/
	//get username and password
	log.Println("METHOD: ", r.Method)

	d.User.Username = r.FormValue("username")
	password := r.FormValue("password")

	if r.Method == "POST" {
		submit := r.FormValue("submit")
		log.Println("SUBMIT: ", submit)
		log.Println("URL: ", r.URL)
		switch submit {
		case "Login":
			d.Backend = r.FormValue("backend")
			cli := client.NewClient(d.Backend)
			log.Println("CASE LOGIN")
			log.Println("LOGIN backend: ", d.Backend)
			if r.FormValue("token") != "" {
				d.User.Token = r.FormValue("token")
			} else {
				//authentication request
				if err = cli.Login(d.User.Username, password); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					log.Println("login err:", err)
					return
				}
				d.User.Token = cli.User.TokenData.Token
			}
			if d.User.Username != "" {
				d.IsLogin = true
			} else {
				d.IsLogin = false
			}
			log.Println("LOGIN: d.User.Token = ", d.User.Token)

		case "Get Jobs":
			//get user data from User sync.Map
			log.Println("START GET JOBS")
			d.Backend = r.FormValue("back")
			log.Println("GET JOBS backend: ", d.Backend)
			d.User.Token = r.FormValue("token")
			log.Println("GET JOBS: d.User.Token = ", d.User.Token)
			cli := client.NewClient(d.Backend)
			cli.User.TokenData.Token = d.User.Token
			memberid, err = strconv.Atoi(r.FormValue("memberid"))
			if err != nil {
				http.Error(w, "invalid member id", http.StatusUnauthorized)
				return
			}
			log.Println("memberid: ", memberid)
			d.User.MemberID = int32(memberid)

			log.Println("d.User.MemberID = ", d.User.MemberID)

			//get list of batch segment jobs
			joblist, err := cli.GetBatchSegmentJobs(d.User.MemberID)
			if err != nil {
				log.Println("getBatchSegmentJobs err: ", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if len(joblist) > 0 {
				d.IsJobs = true
			}

			d.JobList = getJobList(joblist)

		case "Login with Token":
			d.User.Token = r.FormValue("usertoken")

			d.Backend = "xandr"
			log.Println("LOGIN WITH TOKEN backend: ", d.Backend)
			log.Println("usertoken: ", d.User.Token)
			log.Println("Login with token before: ", d.IsLoginWithToken)

			if d.User.Token != "" {
				d.IsLoginWithToken = true
			} else {
				d.IsLoginWithToken = false
			}
			log.Println("Login with token after: ", d.IsLoginWithToken)

		}
	}
	if err := t.ExecuteTemplate(w, "bsstroubleshooter.html", d); err != nil {
		log.Println("Execute bsstroubleshooter.html ", err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
}
