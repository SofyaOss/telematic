package storage

import (
	"time"
)

type Coordinates struct {
	Latitude  float64
	Longitude float64
}

type Car struct {
	ID     int
	Number int
	Speed  int
	//Coords int
	Coordinates
	Date time.Time
}

type DBInterface interface {
	GetTelematic(date1, date2 int64, car int) ([]Car, error)
	GetLatest(cars []int) ([]Car, error)
}
