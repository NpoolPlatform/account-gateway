//nolint:dupl
package user

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	user1 "github.com/NpoolPlatform/account-gateway/pkg/user"
	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateAccount(ctx context.Context, in *npool.UpdateAccountRequest) (*npool.UpdateAccountResponse, error) {
	handler, err := user1.NewHandler(
		ctx,
		user1.WithID(&in.ID, true),
		user1.WithEntID(&in.EntID, true),
		user1.WithAppID(&in.AppID, true),
		user1.WithUserID(&in.UserID, true),
		user1.WithLabels(in.Labels, false),
		user1.WithMemo(in.Memo, false),
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

func (s *Server) UpdateAppUserAccount(ctx context.Context, in *npool.UpdateAppUserAccountRequest) (*npool.UpdateAppUserAccountResponse, error) {
	handler, err := user1.NewHandler(
		ctx,
		user1.WithID(&in.ID, true),
		user1.WithEntID(&in.EntID, true),
		user1.WithAppID(&in.TargetAppID, true),
		user1.WithUserID(&in.TargetUserID, true),
		user1.WithBlocked(in.Blocked, false),
		user1.WithActive(in.Active, false),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateAppUserAccount",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateAppUserAccountResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.UpdateAccount(ctx)
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
