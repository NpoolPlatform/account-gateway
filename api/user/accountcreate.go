package user

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/user"

	user1 "github.com/NpoolPlatform/account-gateway/pkg/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateAccount(ctx context.Context, in *npool.CreateAccountRequest) (*npool.CreateAccountResponse, error) {
	handler, err := user1.NewHandler(
		ctx,
		user1.WithAppID(&in.AppID, true),
		user1.WithUserID(&in.UserID, true),
		user1.WithCoinTypeID(&in.CoinTypeID, true),
		user1.WithUsedFor(&in.UsedFor, true),
		user1.WithAddress(&in.Address, true),
		user1.WithLabels(in.Labels, false),
		user1.WithAccount(in.Account, false),
		user1.WithAccountType(&in.AccountType, true),
		user1.WithVerificationCode(&in.VerificationCode, true),
		user1.WithMemo(in.Memo, false),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateAccount",
			"In", in,
			"Error", err,
		)
		return &npool.CreateAccountResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateAccount(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateAccount",
			"In", in,
			"Error", err,
		)
		return &npool.CreateAccountResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateAccountResponse{
		Info: info,
	}, nil
}
