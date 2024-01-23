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
	WrongUUID string
	Ok        bool
	ErrMsg    string
	SecNum    int
	ValidMsg  string
	Sections  []string //5 sections with len: 8-4-4-4-12
}
