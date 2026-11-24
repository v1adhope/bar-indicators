package entity

const (
	StatusSuccess = "success"
	StatusError   = "error"
)

// Module is the Waybar custom module JSON output.
// See https://github.com/Alexays/Waybar/wiki/Module:-Custom#return-type.
type Module struct {
	Text    string `json:"text"`
	Tooltip string `json:"tooltip"`
}
