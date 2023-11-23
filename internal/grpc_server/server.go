package grpc_server

import (
	"context"
	"log"

	pb "practice/internal/grpc"
	"practice/storage/postgres"

	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedGRPCServiceServer
	Grpc   *grpc.Server
	db     *postgres.TelematicDB
	logger *log.Logger
}

type grpcServer Server

func New(db *postgres.TelematicDB) *Server {
	s := Server{
		Grpc: grpc.NewServer(),
		db:   db,
	}
	pb.RegisterGRPCServiceServer(s.Grpc, (*grpcServer)(&s))
	return &s
}

/*
func (s *Server) Open() {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatalf("could not listen to port: %s", err)
	}
	s.lis = lis
	if err := s.Grpc.Serve(lis); err != nil {
		s.logger.Println("gRPC server returned:", err.Error())
	}
}
*/

func (g *grpcServer) Close() error {
	g.Grpc.GracefulStop()
	g.logger.Println("gRPC server stopped")
	return nil
}

func (g *grpcServer) GetCarsByDate(ctx context.Context, req *pb.CarsByDateRequest) (*pb.CarsByDateResponse, error) {
	firstDate := req.GetFirstDate()
	lastDate := req.GetLastDate()
	nums := req.GetNums()
	res, err := g.db.GetByDate(firstDate, lastDate, nums, ctx)
	if err != nil {
		return nil, err
	}
	return &pb.CarsByDateResponse{
		Cars: res,
	}, nil
	//for
	//return &pb.GetByDateResponse{
	//	Cars: res,
	//}, nil
}

func (g *grpcServer) GetLastCars(ctx context.Context, req *pb.LastCarsRequest) (*pb.LastCarsResponse, error) {
	nums := req.GetNums()
	res, err := g.db.GetByCarNumber(ctx, nums)
	if err != nil {
		return nil, err
	}
	return &pb.LastCarsResponse{
		Cars: res,
	}, nil
}
