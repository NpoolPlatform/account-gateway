package user

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/user"

	constant "github.com/NpoolPlatform/account-gateway/pkg/message/const"

	user1 "github.com/NpoolPlatform/account-gateway/pkg/user"

	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

func (s *Server) CreateWithdrawAccount(
	ctx context.Context,
	in *npool.CreateWithdrawAccountRequest,
) (
	*npool.CreateWithdrawAccountResponse,
	error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateWithdrawAccount")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("CreateWithdrawAccount", "AppID", in.GetAppID(), "error", err)
		return &npool.CreateWithdrawAccountResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if _, err := uuid.Parse(in.GetUserID()); err != nil {
		logger.Sugar().Errorw("CreateWithdrawAccount", "UserID", in.GetUserID(), "error", err)
		return &npool.CreateWithdrawAccountResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if _, err := uuid.Parse(in.GetCoinTypeID()); err != nil {
		logger.Sugar().Errorw("CreateWithdrawAccount", "CoinTypeID", in.GetCoinTypeID(), "error", err)
		return &npool.CreateWithdrawAccountResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if in.GetAddress() == "" {
		logger.Sugar().Errorw("CreateWithdrawAccount", "Address", in.GetAddress())
		return &npool.CreateWithdrawAccountResponse{}, status.Error(codes.InvalidArgument, "invalid address")
	}

	info, err := user1.CreateWithdrawAccount(
		ctx,
		in.GetAppID(),
		in.GetUserID(),
		in.GetCoinTypeID(),
		in.GetAddress(),
		in.GetLabels(),
	)
	if err != nil {
		logger.Sugar().Errorw("CreateWithdrawAccount", "error", err)
		return &npool.CreateWithdrawAccountResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateWithdrawAccountResponse{
		Info: info,
	}, nil
}
