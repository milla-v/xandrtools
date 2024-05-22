package simulator

import (
	"time"
)

type AuthRequest struct {
	Auth struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"auth"`
}

type AuthResponse struct {
	Response struct {
		Status string `json:"status"`
		Token  string `json:"token"`
	} `json:"response"`
}

type BatchSegmentResponse struct {
	Response struct {
		BatchSegmentUploadJob []BatchSegmentUploadJob `json:"batch_segment_upload_job"`
		Count                 int64                   `json:"count"`
		DbgInfo               struct {
			OutputTerm string        `json:"output_term"`
			Version    string        `json:"version"`
			Warnings   []interface{} `json:"warnings"`
		} `json:"dbg_info"`
		NumElements  int64  `json:"num_elements"`
		StartElement int64  `json:"start_element"`
		Status       string `json:"status"`
	} `json:"response"`
}

type BatchSegmentUploadJob struct {
	CompletedTime       string      `json:"completed_time"`
	CreatedOn           string      `json:"created_on"`
	ErrorCode           interface{} `json:"error_code"`
	ErrorLogLines       string      `json:"error_log_lines"`
	ID                  int64       `json:"id"`
	IsBeamFile          bool        `json:"is_beam_file"`
	JobID               string      `json:"job_id"`
	LastModified        string      `json:"last_modified"`
	MemberID            int64       `json:"member_id"`
	NumInactiveSegment  int64       `json:"num_inactive_segment"`
	NumInvalidFormat    int64       `json:"num_invalid_format"`
	NumInvalidSegment   int64       `json:"num_invalid_segment"`
	NumInvalidTimestamp int64       `json:"num_invalid_timestamp"`
	NumInvalidUser      int64       `json:"num_invalid_user"`
	NumOtherError       int64       `json:"num_other_error"`
	NumPastExpiration   int64       `json:"num_past_expiration"`
	NumUnauthSegment    int64       `json:"num_unauth_segment"`
	NumValid            int64       `json:"num_valid"`
	NumValidUser        int64       `json:"num_valid_user"`
	PercentComplete     int64       `json:"percent_complete"`
	Phase               string      `json:"phase"`
	SegmentLogLines     string      `json:"segment_log_lines"`
	StartTime           string      `json:"start_time"`
	TimeToProcess       string      `json:"time_to_process"`
	UploadedTime        string      `json:"uploaded_time"`
	ValidatedTime       string      `json:"validated_time"`
}

type UserData struct {
	Username  string `json: "username"`
	TokenData struct {
		Token          string    `json: "token"`
		ExpirationTime time.Time `json: "expirationTime"`
		MemberId       string    `json: "memberid"`
	} `json: "tokendata"`
}
