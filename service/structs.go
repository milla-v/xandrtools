package service

const notHex = "Not Hexadecimal"

type idtype struct {
	domain string
	number int
}

type segments struct {
	SegID     string
	SegCode   string
	MemberID  string
	Timestamp string
	Value     string
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
