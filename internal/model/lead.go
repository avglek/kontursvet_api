package model

import "time"

type Lead struct {
	ID           int64      `json:"id"`
	Name         string     `json:"name"`
	PhoneDigital string     `json:"phone_digital"`
	PhoneFormat  string     `json:"phone_format"`
	HomeType     string     `json:"home_type"`
	Location     string     `json:"location"`
	Message      string     `json:"message"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type LeadMessage struct {
	Text         Lead        `json:"text"`
	Attachments  []Attachment `json:"attachments"`
}

type Attachment struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Encoding    string `json:"encoding"`
}