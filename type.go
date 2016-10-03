package main

import "time"

type Subject struct {
	Period     int    `json:"period"`
	Name       string `json:"name"`
	Instructor string `json:"instructor"`
	Reason     string `json:"reason"`
}

type Notification struct {
	Location  string    `json:"location"`
	Subjects  []Subject `json:"subjects"`
	Date      time.Time `json:"date"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Error struct {
	Message string `json:"message"`
}
