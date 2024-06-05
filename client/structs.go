package client

import (
	"encoding/json"
	"strings"
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
		Dbg                   DbgInfo                 `json:"dbg_info"`
		NumElements           int64                   `json:"num_elements"`
		StartElement          int64                   `json:"start_element"`
		Status                string                  `json:"status"`
	} `json:"response"`
}

type BssTimestamp time.Time

func (b BssTimestamp) MarshalJSON() ([]byte, error) {
	s := time.Time(b).UTC().Format("2006-01-02 15:03:04")
	return json.Marshal(s)
}

func (b *BssTimestamp) UnmarshalJSON(bytes []byte) error {
	s := strings.Trim(string(bytes), `"`)
	ts, err := time.Parse("2006-01-02 15:03:04", s)
	if err != nil {
		return err
	}

	*b = BssTimestamp(ts)

	return nil
}

func (b BssTimestamp) String() string {
	return time.Time(b).UTC().Format("2006-01-02 15:03:04Z")
}

type BatchSegmentUploadJob struct {
	CompletedTime       BssTimestamp `json:"completed_time"`
	CreatedOn           BssTimestamp `json:"created_on"`
	ErrorCode           interface{}  `json:"error_code"`
	ErrorLogLines       string       `json:"error_log_lines"`
	ID                  int64        `json:"id"`
	IsBeamFile          bool         `json:"is_beam_file"`
	JobID               string       `json:"job_id"`
	LastModified        BssTimestamp `json:"last_modified"`
	MemberID            int32        `json:"member_id"`
	MatchRate           int          `json:"match_rate"`
	NumInactiveSegment  int64        `json:"num_inactive_segment"`
	NumInvalidFormat    int64        `json:"num_invalid_format"`
	NumInvalidSegment   int64        `json:"num_invalid_segment"`
	NumInvalidTimestamp int64        `json:"num_invalid_timestamp"`
	NumInvalidUser      int64        `json:"num_invalid_user"`
	NumOtherError       int64        `json:"num_other_error"`
	NumPastExpiration   int64        `json:"num_past_expiration"`
	NumUnauthSegment    int64        `json:"num_unauth_segment"`
	NumValid            int64        `json:"num_valid"`
	NumValidUser        int64        `json:"num_valid_user"`
	PercentComplete     int64        `json:"percent_complete"`
	Phase               string       `json:"phase"`
	SegmentLogLines     string       `json:"segment_log_lines"`
	StartTime           BssTimestamp `json:"start_time"`
	TimeToProcess       int64        `json:"time_to_process"`
	UploadedTime        BssTimestamp `json:"uploaded_time"`
	ValidatedTime       BssTimestamp `json:"validated_time"`
}

type DbgInfo struct {
	Instance  string        `json:"instance"`
	DbgTime   int           `json:"time"`
	StartTime time.Time     `json:"start_time"`
	Version   string        `json:"version"`
	TraceID   string        `json:"trace_id"`
	Warnings  []interface{} `json:"warnings"`
}

type UserData struct {
	Username  string `json: "username"`
	TokenData struct {
		Token          string    `json:"token"`
		ExpirationTime time.Time `json:"expirationTime"`
		MemberId       int32     `json:"memberid"`
	} `json: "tokendata"`
}
