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
