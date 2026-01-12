package types

import (
	"encoding/json"
)

type PlayerAction struct {
	UserID      string          `json:"user_id"`
	ActionName  string          `json:"action_name"`
	PayloadData json.RawMessage `json:"payload_data"`
}
