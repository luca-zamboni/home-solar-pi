package db

import "time"

type InverterReading struct {
	Reading   int
	CreatedAt time.Time
}

type HeaterLogs struct {
	Reading   int
	CreatedAt time.Time
	Status    string
}
