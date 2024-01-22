package services

import (
	"context"
	"errors"
	pb "github.com/dzoniops/common/pkg/user"
	"github.com/dzoniops/user-service/auth"
	"github.com/dzoniops/user-service/client"
	"github.com/dzoniops/user-service/db"
	"github.com/dzoniops/user-service/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	pb.UnimplementedUserServiceServer
	ReservationClient   client.ReservationClient
	AccommodationClient client.AccommodationClient
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
		return nil, status.Error(codes.NotFound, "Wrong username or password")
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

func (s *Server) UpdatePassword(
	c context.Context,
	req *pb.PasswordRequest,
) (*emptypb.Empty, error) {
	var user models.User
	if result := db.DB.Where(models.User{Username: req.Username}).First(&user); result.Error != nil {
		return nil, status.Error(codes.NotFound, "User not found")
	}
	if !auth.CheckPasswordHash(req.OldPassword, user.Password) {
		return nil, status.Error(codes.InvalidArgument, "Incorrect old password")
	}
	db.DB.Model(&user).Update("password", auth.HashPassword(req.NewPassword))
	return &emptypb.Empty{}, nil
}

func (s *Server) Update(
	c context.Context,
	req *pb.UserUpdateRequest,
) (*pb.RegisterResponse, error) {
	var user models.User
	if result := db.DB.Where(models.User{Username: req.Username}).First(&user); result.Error != nil {
		return nil, status.Error(codes.NotFound, "User not found")
	}
	db.DB.Model(&user).Updates(models.User{
		ID:            user.ID,
		Email:         req.Email,
		Username:      req.Username,
		Name:          req.Name,
		Surname:       req.Surname,
		PlaceOfLiving: req.PlaceOfLiving,
	})
	return &pb.RegisterResponse{Id: user.ID}, nil
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
	err := s.AccommodationClient.DeleteAccommodationsByHost(c, id)
	if err != nil {
		return nil, err
		//return nil, status.Error(codes.Internal,err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) Validate(_ context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	claims, err := auth.ValidateToken(req.AccessToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	var user models.User

	if result := db.DB.Where(&models.User{Username: claims.Username}).First(&user); result.Error != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	if s.checkRoles(req, claims.Role) != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	return &pb.ValidateResponse{
		UserId: user.ID,
	}, nil
}
func (s *Server) checkRoles(req *pb.ValidateRequest, userRole string) error {

	for _, role := range req.Roles {
		if role == userRole {
			return nil
		}
	}
	return errors.New("user does not have needed role")
}
