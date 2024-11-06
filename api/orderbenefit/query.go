package orderbenefit

import (
	"context"

	orderbenefit1 "github.com/NpoolPlatform/account-gateway/pkg/orderbenefit"
	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/orderbenefit"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetAccounts(ctx context.Context, in *npool.GetAccountsRequest) (*npool.GetAccountsResponse, error) {
	handler, err := orderbenefit1.NewHandler(
		ctx,
		orderbenefit1.WithAppID(&in.AppID, true),
		orderbenefit1.WithUserID(&in.UserID, true),
		orderbenefit1.WithOffset(in.GetOffset()),
		orderbenefit1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAccounts",
			"In", in,
			"Error", err,
		)
		return &npool.GetAccountsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetAccounts(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAccounts",
			"In", in,
			"Error", err,
		)
		return &npool.GetAccountsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAccountsResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetAccount(ctx context.Context, in *npool.GetAccountRequest) (*npool.GetAccountResponse, error) {
	handler, err := orderbenefit1.NewHandler(
		ctx,
		orderbenefit1.WithEntID(&in.EntID, true),
		orderbenefit1.WithAppID(&in.AppID, true),
		orderbenefit1.WithUserID(&in.UserID, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAccount",
			"In", in,
			"Error", err,
		)
		return &npool.GetAccountResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.GetAccount(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAccount",
			"In", in,
			"Error", err,
		)
		return &npool.GetAccountResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAccountResponse{
		Info: info,
	}, nil
}
