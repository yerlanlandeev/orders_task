package reporder

import "time"

type Order struct {
	Room      string    `json:"room"`
	UserEmail string    `json:"user_email"`
	From      time.Time `json:"from"`
	To        time.Time `json:"to"`
}
