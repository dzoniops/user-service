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

//
// func (c *ReservationClient) IsEmptyGuestActiveReservations(
// 	ctx context.Context,
// 	id int64,
// ) (bool, error) {
// 	res, err := c.client.ActiveReservationsGuest(ctx, &pb.IdRequest{Id: id})
// 	if err != nil {
// 		return false, err
// 	}
//
// 	if len(res.Reservations) > 0 {
// 		return false, nil
// 	}
// 	return true, nil
// }
