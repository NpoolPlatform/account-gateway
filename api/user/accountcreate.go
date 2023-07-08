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
		user1.WithAppID(&in.AppID),
		user1.WithUserID(&in.UserID),
		user1.WithCoinTypeID(&in.CoinTypeID),
		user1.WithUsedFor(&in.UsedFor),
		user1.WithAddress(&in.Address),
		user1.WithLabels(in.Labels),
		user1.WithAccount(&in.Account),
		user1.WithAccountType(&in.AccountType),
		user1.WithVerificationCode(&in.VerificationCode),
		user1.WithMemo(in.Memo),
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
