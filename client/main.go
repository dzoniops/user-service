package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/dzoniops/user-service/user"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewUserServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// r, err := c.Register(ctx, &pb.RegisterRequest{
	// 	Name:          "Nikola",
	// 	Surname:       "Petrovic",
	// 	Email:         "nikola123@mail.com",
	// 	Username:      "nikola123",
	// 	Password:      "nikola",
	// 	PlaceOfLiving: "Zrenjanin, Tetovska 43",
	// 	Role:          "",
	// })
	// log.Printf("Greeting: %v", r.GetId())
	// if err != nil {
	// 	log.Fatalf("could not greet: %v", err)
	// }
	r1, err := c.Login(ctx, &pb.LoginRequest{
		Username: "nikola123",
		Password: "nikola",
	})
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("Greeting: %s", r1.GetAccessToken())
	r2, err := c.GetUser(ctx, &pb.IdRequest{Id: 1})
	log.Printf("%v", r2)
}
