package transfer

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	transfer1 "github.com/NpoolPlatform/account-gateway/pkg/transfer"
	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/transfer"
)

func (s *Server) DeleteTransfer(ctx context.Context, in *npool.DeleteTransferRequest) (resp *npool.DeleteTransferResponse, err error) {
	handler, err := transfer1.NewHandler(
		ctx,
		transfer1.WithID(&in.TransferID, true),
		transfer1.WithEntID(&in.EntID, true),
		transfer1.WithAppID(&in.AppID, true),
		transfer1.WithUserID(&in.UserID, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteTransfer",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteTransferResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.DeleteTransfer(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteTransfer",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteTransferResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.DeleteTransferResponse{
		Info: info,
	}, nil
}
