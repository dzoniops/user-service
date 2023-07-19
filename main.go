package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/dzoniops/common/pkg/user"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	reservation "github.com/dzoniops/user-service/client"
	"github.com/dzoniops/user-service/db"
	"github.com/dzoniops/user-service/services"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db.InitDB()

	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("PORT")))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	reservationClient := reservation.InitClient(
		fmt.Sprintf("localhost:%s", os.Getenv("RESERVATION_PORT")),
	)
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &services.Server{reservationClient: *reservationClient})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
