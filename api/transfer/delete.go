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

func (s *Server) DeleteTransfer(
	ctx context.Context,
	in *transfer.DeleteTransferRequest,
) (
	resp *transfer.DeleteTransferResponse,
	err error,
) {
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("DeleteTransfer", "AppID", in.GetAppID(), "error", err)
		return &transfer.DeleteTransferResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if _, err := uuid.Parse(in.GetUserID()); err != nil {
		logger.Sugar().Errorw("DeleteTransfer", "UserID", in.GetUserID(), "error", err)
		return &transfer.DeleteTransferResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if _, err := uuid.Parse(in.GetTransferID()); err != nil {
		logger.Sugar().Errorw("DeleteTransfer", "TransferID", in.GetTransferID(), "error", err)
		return &transfer.DeleteTransferResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := mtransfer.DeleteTransfer(ctx, in.GetTransferID(), in.GetAppID(), in.GetUserID())
	if err != nil {
		logger.Sugar().Errorw("CreateTransfer", "error", err)
		return &transfer.DeleteTransferResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &transfer.DeleteTransferResponse{
		Info: info,
	}, nil
}
