package model

type History struct {
	Id       uint64 `json:",omitempty"`
	User_id  uint64 `json:",omitempty"`
	Space_id uint64 `json:",omitempty"`
	Command  string `json:",omitempty"`
	Result   string `json:",omitempty"`
}
