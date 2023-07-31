package client

import (
	"fmt"

	pb "github.com/dzoniops/common/pkg/accommodation"
	"google.golang.org/grpc"
)

type AccommodationClient struct {
	client pb.AccommodationServiceClient
}

func InitAccommodationClient(url string) *AccommodationClient {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		fmt.Println("Could not connect:", err)
	}
	client := pb.NewAccommodationServiceClient(conn)

	return &AccommodationClient{client: client}
}

