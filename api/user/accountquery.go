package user

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/user"

	constant1 "github.com/NpoolPlatform/account-gateway/pkg/const"
	constant "github.com/NpoolPlatform/account-gateway/pkg/message/const"

	user1 "github.com/NpoolPlatform/account-gateway/pkg/user"

	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

func (s *Server) GetAccounts(
	ctx context.Context,
	in *npool.GetAccountsRequest,
) (
	*npool.GetAccountsResponse,
	error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetAccounts")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("GetAccounts", "AppID", in.GetAppID(), "error", err)
		return &npool.GetAccountsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	limit := int32(constant1.DefaultLimit)
	if in.GetLimit() > 0 {
		limit = in.GetLimit()
	}

	infos, total, err := user1.GetAccounts(ctx, in.GetAppID(), in.GetOffset(), limit)
	if err != nil {
		logger.Sugar().Errorw("GetAccounts", "error", err)
		return &npool.GetAccountsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAccountsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
