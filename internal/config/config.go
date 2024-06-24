package config

import "time"

const (
	UpdateURL = "/update/"
)

const (
	ServerHost     = "http://localhost:8080"
	PollInterval   = 2 * time.Second
	ReportInterval = 10 * time.Second
)
