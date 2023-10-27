package storage

import (
	"practice/internal/generator"
	"time"
)

type Car struct {
	ID     int
	Number int
	Speed  int
	Coords *generator.Coordinates
	Date   time.Time
}

type DBInterface interface {
	GetTelematic(date1, date2 int64, car int) ([]Car, error)
	GetLatest(cars []int) ([]Car, error)
}
