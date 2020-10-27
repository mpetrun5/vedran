package models

import "time"

type Record struct {
	ID        int `storm:"id,increment"`
	NodeId    string
	Status    string
	Timestamp time.Time
}
