package server

import (
	"fmt"
	"log"
	"net"

	proto "github.com/highonsemicolon/experiments/zerofail/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func StartServer(serverPort, mongoURI, dbName string) {
	InitMongoDB(mongoURI, dbName)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", serverPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterRecordServiceServer(grpcServer, &RecordServiceServer{})

	reflection.Register(grpcServer)

	log.Println("gRPC server started on :", serverPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
