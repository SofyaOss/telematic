package storage

import (
	"context"
	"time"

	pb "practice/internal/grpc"
)

type Coordinates struct {
	Latitude  float64
	Longitude float64
}

type Car struct {
	ID     int
	Number int
	Speed  int
	Coordinates
	Date time.Time
}

type DBInterface interface {
	CreateTable(ctx context.Context) error
	DropTable(ctx context.Context) error
	AddData(c *Car, ctx context.Context) error
	GetByDate(d1s, d2s string, nums []int64, ctx context.Context) ([]*pb.Car, error)
	GetByCarNumber(ctx context.Context, carNums []int64) ([]*pb.Car, error)
}
