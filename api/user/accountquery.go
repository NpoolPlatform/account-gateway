//nolint:dupl
package user

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/user"

	user1 "github.com/NpoolPlatform/account-gateway/pkg/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetAccounts(ctx context.Context, in *npool.GetAccountsRequest) (*npool.GetAccountsResponse, error) {
	handler, err := user1.NewHandler(
		ctx,
		user1.WithAppID(&in.AppID),
		user1.WithUserID(&in.UserID),
		user1.WithUsedFor(&in.UsedFor),
		user1.WithOffset(in.GetOffset()),
		user1.WithLimit(in.GetLimit()),
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

func (s *Server) GetAppAccounts(ctx context.Context, in *npool.GetAppAccountsRequest) (*npool.GetAppAccountsResponse, error) {
	handler, err := user1.NewHandler(
		ctx,
		user1.WithAppID(&in.AppID),
		user1.WithOffset(in.GetOffset()),
		user1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppAccounts",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppAccountsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetAppAccounts(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppAccounts",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppAccountsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppAccountsResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetNAppAccounts(ctx context.Context, in *npool.GetNAppAccountsRequest) (*npool.GetNAppAccountsResponse, error) {
	handler, err := user1.NewHandler(
		ctx,
		user1.WithAppID(&in.TargetAppID),
		user1.WithOffset(in.GetOffset()),
		user1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetNAppAccounts",
			"In", in,
			"Error", err,
		)
		return &npool.GetNAppAccountsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetAppAccounts(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetNAppAccounts",
			"In", in,
			"Error", err,
		)
		return &npool.GetNAppAccountsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetNAppAccountsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
