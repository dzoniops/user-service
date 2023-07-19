package services

import (
	"context"

	pb "github.com/dzoniops/common/pkg/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/dzoniops/user-service/auth"
	reservation "github.com/dzoniops/user-service/client"
	"github.com/dzoniops/user-service/db"
	"github.com/dzoniops/user-service/models"
)

type Server struct {
	pb.UnimplementedUserServiceServer
	ReservationClient reservation.ReservationClient
}

func (s *Server) Register(
	c context.Context,
	req *pb.RegisterRequest,
) (*pb.RegisterResponse, error) {
	var user models.User
	if result := db.DB.Where(&models.User{Email: req.Email}).First(&user); result.Error == nil {
		return nil, status.Error(codes.AlreadyExists,
			"Email already used")
	}

	if result := db.DB.Where(&models.User{Username: req.Username}).First(&user); result.Error == nil {
		return nil, status.Error(codes.AlreadyExists,
			"Username already used")
	}
	user = models.User{
		Email:         req.Email,
		Username:      req.Username,
		Password:      auth.HashPassword(req.Password),
		Name:          req.Name,
		Surname:       req.Surname,
		PlaceOfLiving: req.PlaceOfLiving,
		Role:          req.Role,
	}

	db.DB.Create(&user)
	return &pb.RegisterResponse{Id: int64(user.ID)}, status.New(codes.OK, "").Err()
}

func (s *Server) Login(c context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	var user models.User

	if result := db.DB.Where(models.User{Username: req.Username}).First(&user); result.Error != nil {
		return nil, status.Error(codes.NotFound, "Wrong username")
	}
	if !auth.CheckPasswordHash(req.Password, user.Password) {
		return nil, status.Error(codes.NotFound, "Wrong username or password")
	}
	accessToken, err := auth.GenerateToken(user)
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}
	return &pb.LoginResponse{AccessToken: accessToken}, nil
}

func (s *Server) GetUser(c context.Context, req *pb.IdRequest) (*pb.UserResponse, error) {
	var user models.User
	if result := db.DB.Where(models.User{ID: req.Id}).First(&user); result.Error != nil {
		return nil, status.Error(codes.NotFound, "User not found")
	}
	data := &pb.UserResponse{
		Name:          user.Name,
		Surname:       user.Surname,
		Email:         user.Email,
		Username:      user.Username,
		PlaceOfLiving: user.PlaceOfLiving,
		Id:            user.ID,
	}
	return data, nil
}

// TODO: password change
func (s *Server) Update(c context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	var user models.User

	if result := db.DB.Where(models.User{Username: req.Username}).First(&user); result.Error != nil {
		return nil, status.Error(codes.NotFound, "User not found")
	}
	user = models.User{
		Email:         req.Email,
		Username:      req.Username,
		Password:      user.Password,
		Name:          req.Name,
		Surname:       req.Surname,
		PlaceOfLiving: req.PlaceOfLiving,
		Role:          req.Role,
	}
	return nil, nil
}

func (s *Server) Delete(c context.Context, req *pb.IdRequest) (*emptypb.Empty, error) {
	var user models.User

	if result := db.DB.Where(models.User{ID: req.Id}).First(&user); result.Error != nil {
		return nil, status.Error(codes.NotFound, "User not found")
	}
	switch user.Role {
	case "GUEST":
		return s.deleteGuest(c, req.Id)
	case "HOST":
		return s.deleteHost(c, req.Id)
	default:
		return nil, status.Error(codes.Unknown, "User role not set to proper one")
	}
}

func (s *Server) deleteGuest(c context.Context, id int64) (*emptypb.Empty, error) {
	isEmpty, err := s.ReservationClient.IsEmptyGuestActiveReservations(c, id)
	if err != nil {
		return &emptypb.Empty{}, err
	}
	if isEmpty {
		db.DB.Delete(&models.User{}, id)
		return &emptypb.Empty{}, nil
	}
	return &emptypb.Empty{}, status.Error(codes.Unavailable, "User has active reservations")
}

func (s *Server) deleteHost(c context.Context, id int64) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
