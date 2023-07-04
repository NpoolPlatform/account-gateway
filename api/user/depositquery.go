//nolint:nolintlint,dupl
package user

import (
	"context"

	user1 "github.com/NpoolPlatform/account-gateway/pkg/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/user"

	"github.com/google/uuid"
)

func (s *Server) GetDepositAccount(ctx context.Context, in *npool.GetDepositAccountRequest) (*npool.GetDepositAccountResponse, error) {
	var err error

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("GetDepositAccount", "AppID", in.GetAppID(), "error", err)
		return &npool.GetDepositAccountResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if _, err := uuid.Parse(in.GetUserID()); err != nil {
		logger.Sugar().Errorw("GetDepositAccount", "UserID", in.GetUserID(), "error", err)
		return &npool.GetDepositAccountResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if _, err := uuid.Parse(in.GetCoinTypeID()); err != nil {
		logger.Sugar().Errorw("GetDepositAccount", "CoinTypeID", in.GetCoinTypeID(), "error", err)
		return &npool.GetDepositAccountResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := user1.GetDepositAccount(ctx, in.GetAppID(), in.GetUserID(), in.GetCoinTypeID())
	if err != nil {
		logger.Sugar().Errorw("GetDepositAccount", "error", err)
		return &npool.GetDepositAccountResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetDepositAccountResponse{
		Info: info,
	}, nil
}

//nolint
func (s *Server) GetDepositAccounts(ctx context.Context, in *npool.GetDepositAccountsRequest) (*npool.GetDepositAccountsResponse, error) { //nolint
	var err error

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("GetDepositAccounts", "AppID", in.GetAppID(), "error", err)
		return &npool.GetDepositAccountsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, n, err := user1.GetDepositAccounts(ctx, in.GetAppID(), in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorw("GetDepositAccounts", "error", err)
		return &npool.GetDepositAccountsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetDepositAccountsResponse{
		Infos: infos,
		Total: n,
	}, nil
}

func (s *Server) GetAppDepositAccounts(ctx context.Context, in *npool.GetAppDepositAccountsRequest) (*npool.GetAppDepositAccountsResponse, error) { //nolint
	var err error

	if _, err := uuid.Parse(in.GetTargetAppID()); err != nil {
		logger.Sugar().Errorw("GetAppDepositAccounts", "TargetAppID", in.GetTargetAppID(), "error", err)
		return &npool.GetAppDepositAccountsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, n, err := user1.GetAppDepositAccounts(ctx, in.GetTargetAppID(), in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorw("GetAppDepositAccounts", "error", err)
		return &npool.GetAppDepositAccountsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppDepositAccountsResponse{
		Infos: infos,
		Total: n,
	}, nil
}
