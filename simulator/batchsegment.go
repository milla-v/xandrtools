package simulator

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"xandrtools/client"
)

func HandleBatchSegment(w http.ResponseWriter, r *http.Request) {
	var err error
	if r.Method != http.MethodGet {
		http.Error(w, "GET", http.StatusMethodNotAllowed)
		return
	}
	/*
		for k, v := range r.Header {
			log.Printf("header: %s=%v", k, v)
		}
	*/
	tokenHeader := r.Header.Get("Authorization")
	token := strings.TrimPrefix(tokenHeader, "Bearer ")
	user, ok := User.Load(token)
	if !ok {
		authenticationFailedError(w)
		return
	}
	u := user.(client.UserData)
	log.Println("u.Token.ExpirationTime: ", u.TokenData.ExpirationTime)
	//check if expiration time exists
	if u.TokenData.ExpirationTime.IsZero() == true {
		http.Error(w, "invalid expiration time: ", http.StatusUnauthorized)
		return
	}
	//check expiration time
	if time.Now().UTC().Before(u.TokenData.ExpirationTime) == false {
		http.Error(w, "invalid expiration time: ", http.StatusUnauthorized)
		return
	}
	log.Println("TIME NOW: ", time.Now().UTC())
	log.Println("EXPIRATION TIME: ", u.TokenData.ExpirationTime)

	s := r.URL.Query().Get("member_id")
	if s == "" {
		noMemberIdError(w)
		return
	}

	log.Println("1. member_id: ", s)

	// var resp BatchSegmentResponse
	var resp client.BatchSegmentResponse
	numJobs := 5
	resp.Response.StartElement = 0
	resp.Response.Count = 1
	resp.Response.Status = "OK"
	log.Println("2. status: ", resp.Response.Status)
	resp.Response.BatchSegmentUploadJob, err = generateBatchSegmentUploadJob(numJobs)
	if err != nil {
		log.Println("generateBatchSegmentUploadJob err", http.StatusUnauthorized)
		return
	}
	log.Println("3. ")
	resp.Response.Dbg, err = generateDbgInfo()
	if err != nil {
		log.Println("generate Dbg-info err", http.StatusUnauthorized)
		return
	}
	log.Println("4. Dbg: ", resp.Response.Dbg.DbgTime)

	buf, err := json.MarshalIndent(resp, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//fmt.Printf("json data: %s\n", buf)
	if _, err := fmt.Fprintln(w, string(buf)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func authenticationFailedError(w io.Writer) {
	var resp client.ErrorResponse

	resp.Response.ErrorId = "NOAUTH"
	resp.Response.Error = "Authentication failed - not logged in"

	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		log.Println("encode error response:", err.Error())
	}

	log.Println("auth failed error")
}

func noMemberIdError(w io.Writer) {
	var resp client.ErrorResponse

	resp.Response.ErrorId = "SYNTAX"
	resp.Response.Error = "no member_id provided"

	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		log.Println("encode error response:", err.Error())
	}

	log.Println("no member_id error")
}

func generateBatchSegmentUploadJob(numJobs int) ([]client.BatchSegmentUploadJob, error) {
	var err error
	var list []client.BatchSegmentUploadJob
	for i := 0; i < numJobs; i++ {
		var u client.BatchSegmentUploadJob
		startTime := time.Now().UTC()
		u.StartTime = client.BssTimestamp(time.Now().UTC())
		u.UploadedTime = client.BssTimestamp(time.Now().UTC().Add(time.Second * 6))
		u.ValidatedTime = client.BssTimestamp(time.Now().UTC().Add(time.Minute * 3))
		completedTime := time.Now().UTC().Add(time.Minute * 1)
		u.CompletedTime = client.BssTimestamp(completedTime)
		u.CreatedOn = client.BssTimestamp(u.StartTime)
		//u.ErrorCode =
		u.ErrorLogLines = "\n\nnum_unauth_segment-4013681496264948522;5013:0,5014:1550"
		u.ID = int64(rand.Int())
		//u.IsBeamFile =
		u.JobID, err = generateToken(20)
		if err != nil {
			log.Println("generate jobId err: ", err)
			return list, err
		}
		u.LastModified = u.CompletedTime
		//u.MemberID = int32(rand.Intn(1000))
		u.NumInactiveSegment = 0
		u.NumInvalidSegment = 0
		//u.NumInvalidTimestamp =
		if i == 0 {
			u.NumInvalidUser = 5000
			u.NumValidUser = 10000
			u.NumInvalidFormat = 2
			u.NumUnauthSegment = 3
		} else {
			u.NumInvalidUser = 500
			u.NumValidUser = 100000
			u.NumInvalidFormat = 0
			u.NumUnauthSegment = 0
		}

		u.NumOtherError = 0
		u.NumPastExpiration = 0
		u.NumValid = 200000
		u.PercentComplete = 100
		u.Phase = "completed"
		u.SegmentLogLines = "\n5010:100000\n5011:50000\n5012:50000"
		// TimeToProcess in Nanosecond
		u.TimeToProcess = int64(completedTime.Sub(startTime))
		list = append(list, u)
	}

	return list, err
}

func generateDbgInfo() (client.DbgInfo, error) {
	var err error
	var dbg client.DbgInfo
	dbg.Instance = "authentication-api-production-8664bd4765-btqsz"
	dbg.DbgTime = 0
	dbg.StartTime = time.Now().UTC()
	dbg.Version = "0.0.0"
	dbg.TraceID, err = generateToken(10)
	if err != nil {
		log.Println("generate trace_id err: ", err)
		return dbg, err
	}
	return dbg, err
}
