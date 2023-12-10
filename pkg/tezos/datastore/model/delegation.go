package model

import "time"

// Delegation represents a delegation model in our datastore.
type Delegation struct {
	Timestamp time.Time
	Amount    int64  `json:"amount"`
	Delegator string `json:"delegator"`
	Block     string `json:"block"`
}
