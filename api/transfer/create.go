package transfer

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	transfer1 "github.com/NpoolPlatform/account-gateway/pkg/transfer"
	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/transfer"
)

func (s *Server) CreateTransfer(ctx context.Context, in *npool.CreateTransferRequest) (resp *npool.CreateTransferResponse, err error) {
	handler, err := transfer1.NewHandler(
		ctx,
		transfer1.WithAppID(&in.AppID, true),
		transfer1.WithUserID(&in.UserID, true),
		transfer1.WithAccount(in.Account, false),
		transfer1.WithAccountType(&in.AccountType, true),
		transfer1.WithVerificationCode(&in.VerificationCode, true),
		transfer1.WithTargetAccount(&in.TargetAccount, true),
		transfer1.WithTargetAccountType(&in.TargetAccountType, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateTransfer",
			"In", in,
			"Error", err,
		)
		return &npool.CreateTransferResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateTransfer(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateTransfer",
			"In", in,
			"Error", err,
		)
		return &npool.CreateTransferResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateTransferResponse{
		Info: info,
	}, nil
}
