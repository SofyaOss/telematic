package grpc_server

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	pb "practice/internal/grpc"
	"practice/storage/postgres"
)

//type Server struct {
//	pb.UnimplementedGRPCServiceServer
//}

type Server struct {
	pb.UnimplementedGRPCServiceServer
	grpc   *grpc.Server
	db     *postgres.TelematicDB
	ln     net.Listener
	addr   string
	logger *log.Logger
}

type grpcServer Server

func New(addr string, db *postgres.TelematicDB) *Server {
	s := Server{
		grpc: grpc.NewServer(),
		db:   db,
		addr: addr,
	}
	pb.RegisterGRPCServiceServer(s.grpc, (*grpcServer)(&s))
	return &s
}

func (s *Server) Open() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	s.ln = ln
	go func() {
		err := s.grpc.Serve(s.ln)
		if err != nil {
			s.logger.Println("gRPC server returned:", err.Error())
		}
	}()
	return nil
}

func (g *grpcServer) Close() error {
	g.grpc.GracefulStop()
	g.logger.Println("gRPC server stopped")
	return nil
}

func (g *grpcServer) GetByDate(ctx context.Context, req *pb.GetByDateRequest) (*pb.GetByDateResponse, error) {
	firstDate := req.GetFirstDate()
	lastDate := req.GetLastDate()
	res, err := g.db.GetByDate(firstDate, lastDate)
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
	return nil, nil
}
