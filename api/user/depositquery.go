//nolint:nolintlint,dupl
package user

import (
	"context"

	user1 "github.com/NpoolPlatform/account-gateway/pkg/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/user"
)

func (s *Server) GetDepositAccount(ctx context.Context, in *npool.GetDepositAccountRequest) (*npool.GetDepositAccountResponse, error) {
	handler, err := user1.NewHandler(
		ctx,
		user1.WithAppID(&in.AppID, true),
		user1.WithUserID(&in.UserID, true),
		user1.WithCoinTypeID(&in.CoinTypeID, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetDepositAccount",
			"In", in,
			"Error", err,
		)
		return &npool.GetDepositAccountResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.GetDepositAccount(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetDepositAccount",
			"In", in,
			"Error", err,
		)
		return &npool.GetDepositAccountResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetDepositAccountResponse{
		Info: info,
	}, nil
}

//nolint
func (s *Server) GetDepositAccounts(ctx context.Context, in *npool.GetDepositAccountsRequest) (*npool.GetDepositAccountsResponse, error) { //nolint
	handler, err := user1.NewHandler(
		ctx,
		user1.WithAppID(&in.AppID, true),
		user1.WithOffset(in.GetOffset()),
		user1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetDepositAccounts",
			"In", in,
			"Error", err,
		)
		return &npool.GetDepositAccountsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetDepositAccounts(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetDepositAccounts",
			"In", in,
			"Error", err,
		)
		return &npool.GetDepositAccountsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetDepositAccountsResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetAppDepositAccounts(ctx context.Context, in *npool.GetAppDepositAccountsRequest) (*npool.GetAppDepositAccountsResponse, error) { //nolint
	handler, err := user1.NewHandler(
		ctx,
		user1.WithAppID(&in.TargetAppID, true),
		user1.WithOffset(in.GetOffset()),
		user1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppDepositAccounts",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppDepositAccountsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetDepositAccounts(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppDepositAccounts",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppDepositAccountsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppDepositAccountsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
