package platform

import (
	"context"

	platform1 "github.com/NpoolPlatform/account-gateway/pkg/platform"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/platform"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateAccount(ctx context.Context, in *npool.UpdateAccountRequest) (*npool.UpdateAccountResponse, error) {
	handler, err := platform1.NewHandler(
		ctx,
		platform1.WithID(&in.ID),
		platform1.WithActive(in.Active),
		platform1.WithLocked(in.Locked),
		platform1.WithBlocked(in.Blocked),
		platform1.WithBackup(in.Backup),
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
