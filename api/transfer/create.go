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
		transfer1.WithAppID(&in.AppID),
		transfer1.WithUserID(&in.UserID),
		transfer1.WithAccount(&in.Account),
		transfer1.WithAccountType(&in.AccountType),
		transfer1.WithVerificationCode(&in.VerificationCode),
		transfer1.WithTargetAccount(&in.TargetAccount),
		transfer1.WithTargetAccountType(&in.TargetAccountType),
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
