//nolint:dupl
package user

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/user"

	user1 "github.com/NpoolPlatform/account-gateway/pkg/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateAccount(
	ctx context.Context,
	in *npool.UpdateAccountRequest,
) (
	*npool.UpdateAccountResponse,
	error,
) {
	handler, err := user1.NewHandler(
		ctx,
		user1.WithID(&in.ID),
		user1.WithAppID(&in.AppID),
		user1.WithUserID(&in.UserID),
		user1.WithLabels(&in.Labels),
		user1.WithMemo(in.Memo),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateAccount",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateAccountResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.UpdateAccount(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateAccount",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateAccountResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateAccountResponse{
		Info: info,
	}, nil
}

//nolint
func (s *Server) UpdateAppUserAccount(
	ctx context.Context,
	in *npool.UpdateAppUserAccountRequest,
) (
	*npool.UpdateAppUserAccountResponse,
	error,
) {
	handler, err := user1.NewHandler(
		ctx,
		user1.WithID(&in.ID),
		user1.WithAppID(&in.TargetAppID),
		user1.WithUserID(&in.TargetUserID),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateAppUserAccount",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateAppUserAccountResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.UpdateAppUserAccount(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateAppUserAccount",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateAppUserAccountResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateAppUserAccountResponse{
		Info: info,
	}, nil
}
