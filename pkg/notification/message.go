package notification

import "encoding/json"

type Type string

const (
	TypePromotion Type = "promotion"
	TypeDiscount  Type = "discount"
)

type Message struct {
	Type    Type            `json:"type"`
	Title   string          `json:"title"`
	Body    string          `json:"body"`
	Payload json.RawMessage `json:"payload,omitempty"`
}
