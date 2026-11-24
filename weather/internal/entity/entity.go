package entity

import "fmt"

type Status int

const (
	StatusSuccess Status = iota
	StatusError
)

func (s Status) String() string {
	switch s {
	case StatusSuccess:
		return "success"
	case StatusError:
		return "error"
	default:
		return fmt.Sprintf("Status(%d)", s)
	}
}

// Module is the Waybar custom module JSON output.
// See https://github.com/Alexays/Waybar/wiki/Module:-Custom#return-type.
type Module struct {
	Text    string `json:"text"`
	Tooltip string `json:"tooltip"`
}
