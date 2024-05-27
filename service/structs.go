package service

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
