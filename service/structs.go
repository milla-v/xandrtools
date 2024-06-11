package service

import (
	"xandrtools/client"
)

const notHex = "Not Hexadecimal"

type idtype struct {
	domain string
	number int
}

type separators struct {
	Sep1 string
	Sep2 string
	Sep3 string
	Sep4 string
	Sep5 string
}

type xandr struct {
	WrongUserID string
	Ok          bool
	ErrList     []string
	UserID      int64
	ValidMsg    string
}

type uuid struct {
	UUID          string
	ErrMsg        string
	ErrSecNum     int
	Sections      []string
	SectionsExist bool
	Ok            bool
}

type Vcs struct {
	RevisionFull  string
	RevisionShort string
	Modified      bool
}

// BSUJ: Batch Segment Upload Job
type BSUJ struct {
	client.BatchSegmentUploadJob
	BSUJerror
}

// BSUJ: Batch Segment Upload Job
type BSUJerror struct {
	CompletedTimeErr       string
	CreatedOnErr           string
	ErrorCodeErr           string
	ErrorLogLinesErr       string
	ErrIDErr               string
	ErrIsBeamFileErr       string
	ErrJobIDErr            string
	LastModifiedErr        string
	MemberIDErr            string
	MatchRateErr           string
	NumInactiveSegmentErr  string
	NumInvalidFormatErr    string
	NumInvalidSegmentErr   string
	NumInvalidTimestampErr string
	NumInvalidUserErr      string
	NumOtherErrorErr       string
	NumPastExpirationErr   string
	NumUnauthSegmentErr    string
	NumValidErr            string
	NumValidUserErr        string
	PercentCompleteErr     string
	PhaseErr               string
	SegmentLogLinesErr     string
	StartTimeErr           string
	TimeToProcessErr       string
	UploadedTimeErr        string
	ValidatedTimeErr       string
}
