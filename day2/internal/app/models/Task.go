package models

import "time"

type Task struct {
	ID   int64     `json:"id"`
	Text string    `json:"text"`
	Tags []string  `json:"tags"`
	Due  time.Time `json:"due"`
}
