package grpc_server

import (
	"context"
	"google.golang.org/grpc"
	"log"
	pb "practice/internal/grpc"
	"practice/storage/postgres"
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

func (g *grpcServer) GetByDate(ctx context.Context, req *pb.GetByDateRequest) (*pb.GetByDateResponse, error) {
	firstDate := req.GetFirstDate()
	lastDate := req.GetLastDate()
	nums := req.GetNums()
	res, err := g.db.GetByDate(firstDate, lastDate, nums)
	if err != nil {
		return nil, err
	}
	return &pb.GetByDateResponse{
		Cars: res,
	}, nil
	//for
	//return &pb.GetByDateResponse{
	//	Cars: res,
	//}, nil
}

func (g *grpcServer) GetLast(ctx context.Context, req *pb.GetLastRequest) (*pb.GetLastResponse, error) {
	nums := req.GetNums()
	res, err := g.db.GetByCarNumber(nums)
	if err != nil {
		return nil, err
	}
	return &pb.GetLastResponse{
		Cars: res,
	}, nil
}
