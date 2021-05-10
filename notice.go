package main

import "time"

type Notification struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	ViewCount int       `json:"view_count"`
}

type NotificationData struct {
	Success bool `json:"success"`
	Data    struct {
		TotalCount int `json:"total_count"`
		TotalPages int `json:"total_pages"`
		List       []Notification `json:"list"`
		FixedNotices []Notification `json:"fixed_notices"`
	} `json:"data"`
}