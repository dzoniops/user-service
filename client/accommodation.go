package client

import (
	"context"
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

func (c *AccommodationClient) DeleteAccommodationsByHost(ctx context.Context, hostId int64) error {
	_, err := c.client.DeleteByHost(ctx, &pb.IdRequest{Id: hostId})
	if err != nil {
		return err
	}
	return nil
}
