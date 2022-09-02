package transfer

import (
	"context"

	constant "github.com/NpoolPlatform/account-gateway/pkg/message/const"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool/account/gw/v1/transfer"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
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
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "DeleteTransfer")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

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
