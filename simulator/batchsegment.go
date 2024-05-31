package simulator

import (
	"encoding/json"
	"fmt"
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

	for k, v := range r.Header {
		log.Printf("header: %s=%v", k, v)
	}

	tokenHeader := r.Header.Get("Authorization")
	token := strings.TrimPrefix(tokenHeader, "Bearer ")

	log.Println("token:", token)

	user, ok := User.Load(token)
	if !ok {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	log.Printf("user: %+v", user)
	log.Println("token expiration time: ", user)

	u := user.(UserData)
	log.Println("u.TokenData. ExpirationTime: ", u.TokenData.ExpirationTime)

	//check if expiration time exists
	if u.TokenData.ExpirationTime.IsZero() == true {
		http.Error(w, "invalid expiration time: ", http.StatusUnauthorized)
		return
	}
	//check expiration time
	if time.Now().Before(u.TokenData.ExpirationTime) == false {
		http.Error(w, "invalid expiration time: ", http.StatusUnauthorized)
		return
	}

	// var resp BatchSegmentResponse
	var resp client.BatchSegmentResponse
	numJobs := 5
	resp.Response.StartElement = 0
	resp.Response.Count = 1
	resp.Response.Status = "OK"
	resp.Response.BatchSegmentUploadJob, err = generateBatchSegmentUploadJob(numJobs)
	if err != nil {
		log.Println("generateBatchSegmentUploadJob err", http.StatusUnauthorized)
		return
	}
	resp.Response.Dbg, err = generateDbgInfo()
	if err != nil {
		log.Println("generate Dbg-info err", http.StatusUnauthorized)
		return
	}

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

func generateBatchSegmentUploadJob(numJobs int) ([]BatchSegmentUploadJob, error) {
	var err error
	var list []client.BatchSegmentUploadJob
	for i := 0; i < numJobs; i++ {
		var u client.BatchSegmentUploadJob
		startTime := time.Now()
		u.StartTime = bssTimestamp(time.Now())
		u.UploadedTime = bssTimestamp(time.Now().Add(time.Second * 6))
		u.ValidatedTime = bssTimestamp(time.Now().Add(time.Minute * 3))
		completedTime := time.Now().Add(time.Minute * 1)
		u.CompletedTime = bssTimestamp(completedTime)
		u.CreatedOn = bssTimestamp(u.StartTime)
		//u.ErrorCode =
		u.ErrorLogLines = "\n\nnum_unauth_segment-4013681496264948522;5013:0,5014:1550"
		u.ID = int64(rand.Uint64())
		//u.IsBeamFile =
		u.JobID, err = generateToken(20)
		if err != nil {
			log.Println("generate jobId err: ", err)
			return list, err
		}
		u.LastModified = u.CompletedTime
		//math.Abs convert negative random numbers to positive
		u.MemberID = int32(rand.Intn(1000))
		u.NumInactiveSegment = 0
		u.NumInvalidFormat = 0
		u.NumInvalidSegment = 0
		//u.NumInvalidTimestamp =
		u.NumInvalidUser = 0
		u.NumOtherError = 0
		u.NumPastExpiration = 0
		u.NumUnauthSegment = 1
		u.NumValid = 200000
		u.NumValidUser = 100000
		u.PercentComplete = 100

		u.Phase = "completed"
		u.SegmentLogLines = "\n5010:100000\n5011:50000\n5012:50000"
		// TimeToProcess in Nanosecond
		u.TimeToProcess = int64(completedTime.Sub(startTime))
		log.Println("TIME TO PROCESS: ", u.TimeToProcess)
		list = append(list, u)

		log.Println("FOR len uploadJob.TimeToProcess: ", len(list))
	}

	return list, err
}

func generateDbgInfo() (DbgInfo, error) {
	var err error
	var dbg client.DbgInfo
	dbg.Instance = "authentication-api-production-8664bd4765-btqsz"
	dbg.DbgTime = 0
	dbg.StartTime = time.Now()
	dbg.Version = "0.0.0"
	dbg.TraceID, err = generateToken(10)
	if err != nil {
		log.Println("generate trace_id err: ", err)
		return dbg, err
	}
	return dbg, err
}
