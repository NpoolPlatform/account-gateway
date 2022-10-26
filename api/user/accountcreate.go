package user

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/user"
	accountmgrpb "github.com/NpoolPlatform/message/npool/account/mgr/v1/account"

	constant "github.com/NpoolPlatform/account-gateway/pkg/message/const"

	user1 "github.com/NpoolPlatform/account-gateway/pkg/user"

	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

func (s *Server) CreateAccount(
	ctx context.Context,
	in *npool.CreateAccountRequest,
) (
	*npool.CreateAccountResponse,
	error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateAccount")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("CreateAccount", "AppID", in.GetAppID(), "error", err)
		return &npool.CreateAccountResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if _, err := uuid.Parse(in.GetUserID()); err != nil {
		logger.Sugar().Errorw("CreateAccount", "UserID", in.GetUserID(), "error", err)
		return &npool.CreateAccountResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if _, err := uuid.Parse(in.GetCoinTypeID()); err != nil {
		logger.Sugar().Errorw("CreateAccount", "CoinTypeID", in.GetCoinTypeID(), "error", err)
		return &npool.CreateAccountResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if in.GetAddress() == "" {
		logger.Sugar().Errorw("CreateAccount", "Address", in.GetAddress())
		return &npool.CreateAccountResponse{}, status.Error(codes.InvalidArgument, "invalid address")
	}
	switch in.GetUsedFor() {
	case accountmgrpb.AccountUsedFor_UserWithdraw:
	case accountmgrpb.AccountUsedFor_UserDirectBenefit:
	default:
		logger.Sugar().Errorw("CreateAccount", "UsedFor", in.GetUsedFor())
		return &npool.CreateAccountResponse{}, status.Error(codes.InvalidArgument, "invalid used for")
	}

	info, err := user1.CreateAccount(
		ctx,
		in.GetAppID(),
		in.GetUserID(),
		in.GetCoinTypeID(),
		in.GetUsedFor(),
		in.GetAddress(),
		in.GetLabels(),
	)
	if err != nil {
		logger.Sugar().Errorw("CreateAccount", "error", err)
		return &npool.CreateAccountResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateAccountResponse{
		Info: info,
	}, nil
}
