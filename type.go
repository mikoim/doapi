package main

import "time"

type Class struct {
	Period     int    `json:"period"`
	Name       string `json:"name"`
	Instructor string `json:"instructor"`
	Reason     string `json:"reason"`
}

type Notification struct {
	Location  string    `json:"location"`
	Classes   []Class   `json:"classes"`
	Date      time.Time `json:"date"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Error struct {
	Message string `json:"message"`
}
