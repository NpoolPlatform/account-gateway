//nolint:nolintlint,dupl
package user

import (
	"context"

	commontracer "github.com/NpoolPlatform/account-gateway/pkg/tracer"

	user1 "github.com/NpoolPlatform/account-gateway/pkg/user"

	constant "github.com/NpoolPlatform/account-gateway/pkg/message/const"

	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/user"

	"github.com/google/uuid"
)

func (s *Server) GetDepositAccount(ctx context.Context, in *npool.GetDepositAccountRequest) (*npool.GetDepositAccountResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetDepositAccount")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

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

	span = commontracer.TraceInvoker(span, "user", "user", "GetDepositAccount")

	info, err := user1.GetDepositAccount(ctx, in.GetAppID(), in.GetUserID(), in.GetCoinTypeID())
	if err != nil {
		logger.Sugar().Errorw("GetDepositAccount", "error", err)
		return &npool.GetDepositAccountResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetDepositAccountResponse{
		Info: info,
	}, nil
}