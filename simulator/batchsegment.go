package simulator

import (
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
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

	//var resp BatchSegmentResponse
	uploadJob, err := generateBatchSegmentUploadJob()
	if err != nil {
		log.Println("generateBatchSegmentUploadJob err: ", err)
		return
	}
	log.Println("----------- Generated Batch Segment Upload Job --------------")
	for _, u := range uploadJob {
		log.Println("start_time: ", u.StartTime)
		log.Println("uploaded_time: ", u.UploadedTime)
		log.Println("validated_time: ", u.ValidatedTime)
		log.Println("completed_time: ", u.CompletedTime)
		log.Println("time_to_process: ", u.TimeToProcess)
	}
}

func generateBatchSegmentUploadJob() ([]BatchSegmentUploadJob, error) {
	var err error
	var uploadJob []BatchSegmentUploadJob

	var u BatchSegmentUploadJob
	u.StartTime = time.Now()
	u.UploadedTime = u.StartTime.Add(time.Second * 6)
	u.ValidatedTime = u.UploadedTime.Add(time.Minute * 3)
	u.CompletedTime = u.ValidatedTime.Add(time.Minute * 1)
	u.CreatedOn = u.StartTime
	//u.ErrorCode =
	u.ErrorLogLines = "\n\nnum_unauth_segment-4013681496264948522;5013:0,5014:1550"
	u.ID = int64(rand.Uint64())
	//u.IsBeamFile =
	u.JobID, err = generateToken(20)
	if err != nil {
		log.Println("generate jobId err: ", err)
		return uploadJob, err
	}
	u.LastModified = u.CompletedTime
	u.MemberID = int32(rand.Uint32())
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
	u.TimeToProcess = u.CompletedTime.Sub(u.StartTime)
	uploadJob = append(uploadJob, u)

	log.Println("FOR len uploadJob.TimeToProcess: ", len(uploadJob))
	return uploadJob, err
}
