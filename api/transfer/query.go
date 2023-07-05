//nolint:dupl
package transfer

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	transfer1 "github.com/NpoolPlatform/account-gateway/pkg/transfer"
	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/transfer"
)

func (s *Server) GetTransfers(ctx context.Context, in *npool.GetTransfersRequest) (resp *npool.GetTransfersResponse, err error) {
	handler, err := transfer1.NewHandler(
		ctx,
		transfer1.WithAppID(&in.AppID),
		transfer1.WithUserID(&in.UserID),
		transfer1.WithOffset(in.GetOffset()),
		transfer1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetTransfers",
			"In", in,
			"Error", err,
		)
		return &npool.GetTransfersResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetTransfers(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetTransfers",
			"In", in,
			"Error", err,
		)
		return &npool.GetTransfersResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetTransfersResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetAppTransfers(ctx context.Context, in *npool.GetAppTransfersRequest) (resp *npool.GetAppTransfersResponse, err error) {
	handler, err := transfer1.NewHandler(
		ctx,
		transfer1.WithAppID(&in.AppID),
		transfer1.WithOffset(in.GetOffset()),
		transfer1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppTransfers",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppTransfersResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetTransfers(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppTransfers",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppTransfersResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppTransfersResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetNAppTransfers(ctx context.Context, in *npool.GetNAppTransfersRequest) (resp *npool.GetNAppTransfersResponse, err error) {
	handler, err := transfer1.NewHandler(
		ctx,
		transfer1.WithAppID(&in.TargetAppID),
		transfer1.WithOffset(in.GetOffset()),
		transfer1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetNAppTransfers",
			"In", in,
			"Error", err,
		)
		return &npool.GetNAppTransfersResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetTransfers(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetNAppTransfers",
			"In", in,
			"Error", err,
		)
		return &npool.GetNAppTransfersResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetNAppTransfersResponse{
		Infos: infos,
		Total: total,
	}, nil
}
