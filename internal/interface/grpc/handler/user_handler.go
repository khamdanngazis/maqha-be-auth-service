package handler

import (
	"context"
	"maqhaa/auth_service/internal/app/service"
	pb "maqhaa/auth_service/internal/interface/grpc/model" // Update with your actual package name
)

type UserHandler struct {
	userService service.AuthService
}

func NewUserGRPCHandler(UserService service.AuthService) *UserHandler {
	return &UserHandler{
		userService: UserService,
	}
}
func (h *UserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	User, appError := h.userService.Authorize(ctx, req.Token)
	var response *pb.GetUserResponse

	if appError.Code != service.SuccessError {
		response = &pb.GetUserResponse{
			Code:    int32(appError.Code),
			Message: appError.Message,
			Data:    nil,
		}
		return response, nil
	}

	response = &pb.GetUserResponse{
		Code:    int32(appError.Code),
		Message: appError.Message,
		Data: &pb.UserData{
			Id:       uint32(User.ID),
			ClientId: uint32(User.ClientID),
			IsLogin:  User.IsLogin,
			IsAdmin:  User.IsAdmin,
		},
	}
	return response, nil
}
