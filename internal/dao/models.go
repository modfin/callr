package dao

import "time"

type Person struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type Incident struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Status    string    `json:"status"`

	OnCallIndex int `json:"on_call_index"`

	CallId   string    `json:"call_id"`
	LastCall time.Time `json:"last_call"`

	Responsible *Person `json:"responsible"`

	Declined []Person `json:"declined,omitempty"`
}

type Log struct {
	CreatedAt   time.Time           `json:"created_at"`
	ContentType string              `json:"content_type"`
	Params      map[string][]string `json:"params"`
	Body        string              `json:"body"`
}
