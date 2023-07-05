package user

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/user"

	user1 "github.com/NpoolPlatform/account-gateway/pkg/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) DeleteAccount(
	ctx context.Context,
	in *npool.DeleteAccountRequest,
) (
	*npool.DeleteAccountResponse,
	error,
) {
	handler, err := user1.NewHandler(
		ctx,
		user1.WithID(&in.ID),
		user1.WithAppID(&in.AppID),
		user1.WithUserID(&in.UserID),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteAccount",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteAccountResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.DeleteAccount(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteAccount",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteAccountResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.DeleteAccountResponse{
		Info: info,
	}, nil
}
