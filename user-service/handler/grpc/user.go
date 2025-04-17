package grpc

import (
	"context"
	"proto/user/userpb"
	"user-service/domain"
	"user-service/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	userpb.UnimplementedUserServiceServer
	uc usecase.UserUsecase
}

func NewUserHandler(uc usecase.UserUsecase) *UserHandler {
	return &UserHandler{uc: uc}
}

func (h *UserHandler) RegisterUser(ctx context.Context, req *userpb.RegisterRequest) (*userpb.UserResponse, error) {
	user := &domain.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	err := h.uc.Register(user)
	if err != nil {
		return &userpb.UserResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &userpb.UserResponse{
		Success: true,
		Message: "User registered successfully",
		User: &userpb.User{
			Id:       int32(user.ID),
			Username: user.Username,
			Email:    user.Email,
		},
	}, nil
}

func (h *UserHandler) AuthenticateUser(ctx context.Context, req *userpb.AuthRequest) (*userpb.UserResponse, error) {
	user, err := h.uc.Authenticate(req.Email, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
	}

	return &userpb.UserResponse{
		Success: true,
		Message: "Authenticated",
		User: &userpb.User{
			Id:       int32(user.ID),
			Username: user.Username,
			Email:    user.Email,
		},
	}, nil
}

func (h *UserHandler) GetUserProfile(ctx context.Context, req *userpb.UserID) (*userpb.User, error) {
	user, err := h.uc.GetProfile(int(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	return &userpb.User{
		Id:       int32(user.ID),
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (h *UserHandler) UpdateUserProfile(ctx context.Context, req *userpb.UpdateRequest) (*userpb.UserResponse, error) {
	user := &domain.User{
		ID:       int(req.Id),
		Username: req.Username,
		Email:    req.Email,
	}

	err := h.uc.UpdateProfile(user)
	if err != nil {
		return &userpb.UserResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &userpb.UserResponse{
		Success: true,
		Message: "Profile updated",
		User: &userpb.User{
			Id:       req.Id,
			Username: req.Username,
			Email:    req.Email,
		},
	}, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, req *userpb.UserID) (*userpb.UserResponse, error) {
	err := h.uc.DeleteUser(int(req.Id))
	if err != nil {
		return &userpb.UserResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}
	return &userpb.UserResponse{
		Success: true,
		Message: "User deleted",
	}, nil
}

func (h *UserHandler) ListUsers(ctx context.Context, _ *userpb.Empty) (*userpb.UserList, error) {
	users, err := h.uc.ListUsers()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not list users")
	}

	var pbUsers []*userpb.User
	for _, u := range users {
		pbUsers = append(pbUsers, &userpb.User{
			Id:       int32(u.ID),
			Username: u.Username,
			Email:    u.Email,
		})
	}

	return &userpb.UserList{Users: pbUsers}, nil
}
