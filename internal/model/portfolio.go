package model

import "time"

type PortfolioCase struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Case        string    `json:"case"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Task        string    `json:"task"`
	Works       []string  `json:"works"`
	Location    string    `json:"location"`
	Term        string    `json:"term"`
	Team        string    `json:"team"`
	Period      string    `json:"period"`
	Features    string    `json:"features"`
	Meta        []string  `json:"meta"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PortfolioPhoto struct {
	ID           int64     `json:"id"`
	CaseID       int64     `json:"case_id"`
	GalleryKey   int       `json:"gallery_key"`
	PhotoPath    string    `json:"photo_path"`
	Alt          string    `json:"alt"`
	Caption      string    `json:"caption"`
	CreatedAt    time.Time `json:"created_at"`
}