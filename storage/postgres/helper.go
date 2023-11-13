package postgres

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	pb "practice/internal/grpc"
	"practice/storage"
)

func CarPostgresToProto(pgCar *storage.Car) *pb.Car {
	return &pb.Car{
		Id:        int64(pgCar.ID),
		Number:    int64(pgCar.Number),
		Speed:     int64(pgCar.Speed),
		Latitude:  float32(pgCar.Coordinates.Latitude),
		Longitude: float32(pgCar.Coordinates.Longitude),
		Date:      timestamppb.New(pgCar.Date),
	}
}
