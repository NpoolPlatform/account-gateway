package goodbenefit

import (
	"context"

	goodbenefit1 "github.com/NpoolPlatform/account-gateway/pkg/goodbenefit"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/goodbenefit"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateAccount(ctx context.Context, in *npool.UpdateAccountRequest) (*npool.UpdateAccountResponse, error) {
	handler, err := goodbenefit1.NewHandler(
		ctx,
		goodbenefit1.WithID(&in.ID, true),
		goodbenefit1.WithEntID(&in.EntID, true),
		goodbenefit1.WithActive(in.Active, false),
		goodbenefit1.WithLocked(in.Locked, false),
		goodbenefit1.WithBlocked(in.Blocked, false),
		goodbenefit1.WithBackup(in.Backup, false),
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
