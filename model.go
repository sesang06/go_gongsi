package main

import "time"

type NoticePost struct {
	Assets    string    `json:"assets"`
	ID        int       `json:"id"`
	StartDate time.Time `json:"start_date"`
	Text      string    `json:"text"`
	URL       string    `json:"url"`
}

type NoticeObject struct {
	Data struct {
		More   bool `json:"more"`
		Offset int  `json:"offset"`
		Posts  []NoticePost `json:"posts"`
	} `json:"data"`
	Success bool   `json:"success"`
	Time    string `json:"time"`
}
