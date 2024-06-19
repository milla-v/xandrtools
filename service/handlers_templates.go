package service

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/milla-v/xandr/bss/xgen"

	"xandrtools/client"
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

		Auth client.AuthRequest
		User client.UserData

		IsJobs         bool
		Token          string
		Backend        string
		ExpirationTime time.Time
		JobList        []WebsiteBSUJ
	}
	var d data
	var err error

	d.XandrVersion = Version
	d.VCS.RevisionFull = VcsInfo.RevisionFull
	d.VCS.RevisionShort = VcsInfo.RevisionShort
	d.VCS.Modified = VcsInfo.Modified

	//get username and password
	log.Println("METHOD: ", r.Method)

	d.Auth.Auth.Username = r.FormValue("username")
	d.Auth.Auth.Password = r.FormValue("password")
	d.Backend = r.FormValue("backend")
	if r.Method == "POST" {
		submit := r.FormValue("submit")
		log.Println("SUBMIT: ", submit)
		switch submit {
		case "Login":
			log.Println("CASE LOGIN")
			cli := client.NewClient(d.Backend)

			if r.FormValue("token") != "" {
				cli.User.TokenData.Token = r.FormValue("token")
			} else {
				//authentication request
				d.Auth.Auth.Username = r.FormValue("username")
				d.Auth.Auth.Password = r.FormValue("password")

				if err = cli.Login(r.FormValue("username"), r.FormValue("password")); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				log.Println("client token: ", cli.User.TokenData.Token, "token expiration time: ", cli.User.TokenData.ExpirationTime)
				d.User = cli.User
				d.User.Username = d.Auth.Auth.Username
				cli.User.Username = d.Auth.Auth.Username

			}

			d.Token = cli.User.TokenData.Token

		case "Get Jobs":
			log.Println("CASE GET JOBS")
			cli := client.NewClient(d.Backend)

			//get user data from User sync.Map
			cli.User.TokenData.Token = r.FormValue("token")
			memberid, err := strconv.Atoi(r.FormValue("memberid"))
			if err != nil {
				http.Error(w, "invalid member id", http.StatusUnauthorized)
				return
			}
			cli.User.TokenData.MemberId = int32(memberid)

			user, ok := simulator.User.Load(cli.User.TokenData.Token)
			if !ok {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			d.User = user.(client.UserData)

			//get list of batch segment jobs
			jobs, err := cli.GetBatchSegmentJobs(cli.User.TokenData.MemberId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if len(jobs) > 0 {
				d.IsJobs = true
			} else {
				d.IsJobs = false
			}

			d.Token = d.User.TokenData.Token
			d.User.Username = r.FormValue("username")
			d.User.TokenData.MemberId = cli.User.TokenData.MemberId
			d.JobList = make([]WebsiteBSUJ, len(jobs))
			for i, job := range jobs {
				d.JobList[i].BatchSegmentUploadJob = job
				d.JobList[i].BatchSegmentUploadJob.MatchRate = int(d.JobList[i].BatchSegmentUploadJob.NumValidUser * 100 / (d.JobList[i].BatchSegmentUploadJob.NumValidUser + d.JobList[i].BatchSegmentUploadJob.NumInvalidUser))
				if d.JobList[i].BatchSegmentUploadJob.MatchRate < 71 {
					d.JobList[i].BSUJerror.MatchRateErr = "Low match rate."
					d.JobList[i].JobErrors = append(d.JobList[i].JobErrors, d.JobList[i].BSUJerror.MatchRateErr)
				}
				if d.JobList[i].BatchSegmentUploadJob.ErrorLogLines != "" && d.JobList[i].BatchSegmentUploadJob.MatchRate < 71 {
					d.JobList[i].BSUJerror.ErrorLogLinesErr = "Remove invalid segments."
					d.JobList[i].JobErrors = append(d.JobList[i].JobErrors, d.JobList[i].BSUJerror.ErrorLogLinesErr)
				}
				if d.JobList[i].NumInvalidFormat > 0 {
					d.JobList[i].BSUJerror.NumInvalidFormatErr = "Fix invalid format."
					d.JobList[i].JobErrors = append(d.JobList[i].JobErrors, d.JobList[i].BSUJerror.NumInvalidFormatErr)
				}
				if d.JobList[i].NumUnauthSegment > 0 {
					d.JobList[i].NumUnauthSegmentErr = "Remove num_unauth_segment or verify that segment is active using apixandr.com/segment API call."
					d.JobList[i].JobErrors = append(d.JobList[i].JobErrors, d.JobList[i].NumUnauthSegmentErr)
				}
			}

		}
	}
	if err := t.ExecuteTemplate(w, "bsstroubleshooter.html", d); err != nil {
		log.Println("Execute bsstroubleshooter.html ", err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
}

func handleBssTroubleShooterTest(w http.ResponseWriter, r *http.Request) {
	log.Println("start bss trouble shooter")
	type data struct {
		XandrVersion string
		VCS          Vcs

		User client.UserData

		Backend        string
		ExpirationTime time.Time
		JobList        []WebsiteBSUJ
		IsJobs         bool
	}
	var d data
	var err error
	var password string

	d.XandrVersion = Version
	d.VCS.RevisionFull = VcsInfo.RevisionFull
	d.VCS.RevisionShort = VcsInfo.RevisionShort
	d.VCS.Modified = VcsInfo.Modified

	//get username and password
	log.Println("METHOD: ", r.Method)

	d.User.Username = r.FormValue("username")
	password = r.FormValue("password")
	d.Backend = r.FormValue("backend")
	if r.Method == "POST" {
		submit := r.FormValue("submit")
		log.Println("SUBMIT: ", submit)
		switch submit {
		case "Login":
			log.Println("CASE LOGIN")
			cli := client.NewClient(d.Backend)

			if r.FormValue("token") != "" {
				d.Token = r.FormValue("token")
			} else {
				//authentication request
				if err = cli.Login(d.User.Username, password); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				d.User = cli.User
				log.Println("UserName: ", d.User.Username, " | Token: ", d.User.TokenData.Token, " | ExpTime: ", d.User.TokenData.ExpirationTime, " | MemberID: ", d.User.TokenData.MemberId)
			}
		case "Get Jobs":
			log.Println("CASE GET JOBS")
			cli := client.NewClient(d.Backend)
			if 

			//get user data from User sync.Map
			cli.User.TokenData.Token = r.FormValue("token")
			memberid, err := strconv.Atoi(r.FormValue("memberid"))
			if err != nil {
				http.Error(w, "invalid member id", http.StatusUnauthorized)
				return
			}
			cli.User.TokenData.MemberId = int32(memberid)

			user, ok := simulator.User.Load(cli.User.TokenData.Token)
			if !ok {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			d.User = user.(client.UserData)

			//get list of batch segment jobs
			jobs, err := cli.GetBatchSegmentJobs(cli.User.TokenData.MemberId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if len(jobs) > 0 {
				d.IsJobs = true
			} else {
				d.IsJobs = false
			}

			d.Token = d.User.TokenData.Token
			d.User.Username = r.FormValue("username")
			d.User.TokenData.MemberId = cli.User.TokenData.MemberId
			d.JobList = make([]WebsiteBSUJ, len(jobs))
			for i, job := range jobs {
				d.JobList[i].BatchSegmentUploadJob = job
				d.JobList[i].BatchSegmentUploadJob.MatchRate = int(d.JobList[i].BatchSegmentUploadJob.NumValidUser * 100 / (d.JobList[i].BatchSegmentUploadJob.NumValidUser + d.JobList[i].BatchSegmentUploadJob.NumInvalidUser))
				if d.JobList[i].BatchSegmentUploadJob.MatchRate < 71 {
					d.JobList[i].BSUJerror.MatchRateErr = "Low match rate."
					d.JobList[i].JobErrors = append(d.JobList[i].JobErrors, d.JobList[i].BSUJerror.MatchRateErr)
				}
				if d.JobList[i].BatchSegmentUploadJob.ErrorLogLines != "" && d.JobList[i].BatchSegmentUploadJob.MatchRate < 71 {
					d.JobList[i].BSUJerror.ErrorLogLinesErr = "Remove invalid segments."
					d.JobList[i].JobErrors = append(d.JobList[i].JobErrors, d.JobList[i].BSUJerror.ErrorLogLinesErr)
				}
				if d.JobList[i].NumInvalidFormat > 0 {
					d.JobList[i].BSUJerror.NumInvalidFormatErr = "Fix invalid format."
					d.JobList[i].JobErrors = append(d.JobList[i].JobErrors, d.JobList[i].BSUJerror.NumInvalidFormatErr)
				}
				if d.JobList[i].NumUnauthSegment > 0 {
					d.JobList[i].NumUnauthSegmentErr = "Remove num_unauth_segment or verify that segment is active using apixandr.com/segment API call."
					d.JobList[i].JobErrors = append(d.JobList[i].JobErrors, d.JobList[i].NumUnauthSegmentErr)
				}
			}

		}
	}
	if err := t.ExecuteTemplate(w, "bsstroubleshootertest.html", d); err != nil {
		log.Println("Execute bsstroubleshooter.html ", err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
}
