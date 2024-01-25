package service

const notHex = "Not Hexadecimal"

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
