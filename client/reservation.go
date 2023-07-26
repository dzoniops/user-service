package client 

import (
	"context"
	"fmt"

	pb "github.com/dzoniops/common/pkg/reservation"
	"google.golang.org/grpc"
)

type ReservationClient struct {
	client pb.ReservationServiceClient
}

func InitReservationClient(url string) *ReservationClient {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		fmt.Println("Could not connect:", err)
	}
	client := pb.NewReservationServiceClient(conn)

	return &ReservationClient{client: client}
}

func (c *ReservationClient) IsEmptyGuestActiveReservations(
	ctx context.Context,
	id int64,
) (bool, error) {
	res, err := c.client.ActiveReservationsGuest(ctx, &pb.IdRequest{Id: id})
	if err != nil {
		return false, err
	}

	if len(res.Reservations) > 0 {
		return false, nil
	}
	return true, nil
}
