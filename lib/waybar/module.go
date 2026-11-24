package waybar

import (
	"encoding/json/v2"
	"fmt"
)

// Module is the Waybar custom module JSON output.
// See https://github.com/Alexays/Waybar/wiki/Module:-Custom#return-type.
type Module struct {
	Text    string `json:"text"`
	Tooltip string `json:"tooltip,omitempty"`
}

func (m *Module) Encode() (string, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("marshal: %w", err)
	}
	return string(data), nil
}
