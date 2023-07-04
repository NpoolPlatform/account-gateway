//nolint:dupl
package transfer

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool/account/gw/v1/transfer"
	"github.com/google/uuid"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	mtransfer "github.com/NpoolPlatform/account-gateway/pkg/transfer"
)

func (s *Server) GetTransfers(ctx context.Context, in *transfer.GetTransfersRequest) (resp *transfer.GetTransfersResponse, err error) {
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("GetTransfers", "AppID", in.GetAppID(), "error", err)
		return &transfer.GetTransfersResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if _, err := uuid.Parse(in.GetUserID()); err != nil {
		logger.Sugar().Errorw("GetTransfers", "UserID", in.GetUserID(), "error", err)
		return &transfer.GetTransfersResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := mtransfer.GetTransfers(ctx, in.GetAppID(), in.GetUserID(), in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorw("GetTransfers", "error", err)
		return &transfer.GetTransfersResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &transfer.GetTransfersResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetAppTransfers(
	ctx context.Context,
	in *transfer.GetAppTransfersRequest,
) (
	resp *transfer.GetAppTransfersResponse,
	err error,
) {
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("GetTransfers", "AppID", in.GetAppID(), "error", err)
		return &transfer.GetAppTransfersResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := mtransfer.GetAppTransfers(ctx, in.GetAppID(), in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorw("GetTransfers", "error", err)
		return &transfer.GetAppTransfersResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &transfer.GetAppTransfersResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetNAppTransfers(
	ctx context.Context,
	in *transfer.GetNAppTransfersRequest,
) (
	resp *transfer.GetNAppTransfersResponse,
	err error,
) {
	if _, err := uuid.Parse(in.GetTargetAppID()); err != nil {
		logger.Sugar().Errorw("GetNAppTransfers", "TargetAppID", in.GetTargetAppID(), "error", err)
		return &transfer.GetNAppTransfersResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := mtransfer.GetAppTransfers(ctx, in.GetTargetAppID(), in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorw("GetNAppTransfers", "error", err)
		return &transfer.GetNAppTransfersResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &transfer.GetNAppTransfersResponse{
		Infos: infos,
		Total: total,
	}, nil
}
