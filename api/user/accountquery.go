//nolint:dupl
package user

//nolint
import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/user"

	constant1 "github.com/NpoolPlatform/account-gateway/pkg/const"
	constant "github.com/NpoolPlatform/account-gateway/pkg/message/const"

	user1 "github.com/NpoolPlatform/account-gateway/pkg/user"

	accountmgrpb "github.com/NpoolPlatform/message/npool/account/mgr/v1/account"

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
	if _, err := uuid.Parse(in.GetUserID()); err != nil {
		logger.Sugar().Errorw("GetAccounts", "UserID", in.GetUserID(), "error", err)
		return &npool.GetAccountsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	switch in.GetUsedFor() {
	case accountmgrpb.AccountUsedFor_UserWithdraw, accountmgrpb.AccountUsedFor_UserDirectBenefit:
	default:
		logger.Sugar().Errorw("GetAccounts", "UsedFor", in.GetUsedFor())
		return &npool.GetAccountsResponse{}, status.Error(codes.InvalidArgument, "invalid usedFor")
	}

	limit := int32(constant1.DefaultLimit)
	if in.GetLimit() > 0 {
		limit = in.GetLimit()
	}

	infos, total, err := user1.GetAccounts(ctx, in.GetAppID(), in.GetUserID(), in.GetUsedFor(), in.GetOffset(), limit)
	if err != nil {
		logger.Sugar().Errorw("GetAccounts", "error", err)
		return &npool.GetAccountsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAccountsResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetAppAccounts(
	ctx context.Context,
	in *npool.GetAppAccountsRequest,
) (
	*npool.GetAppAccountsResponse,
	error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetAppAccounts")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("GetAppAccounts", "AppID", in.GetAppID(), "error", err)
		return &npool.GetAppAccountsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	limit := int32(constant1.DefaultLimit)
	if in.GetLimit() > 0 {
		limit = in.GetLimit()
	}

	infos, total, err := user1.GetAppAccounts(ctx, in.GetAppID(), in.GetOffset(), limit)
	if err != nil {
		logger.Sugar().Errorw("GetAppAccounts", "error", err)
		return &npool.GetAppAccountsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppAccountsResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetNAppAccounts(
	ctx context.Context,
	in *npool.GetNAppAccountsRequest,
) (
	*npool.GetNAppAccountsResponse,
	error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetNAppAccounts")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetTargetAppID()); err != nil {
		logger.Sugar().Errorw("GetNAppAccounts", "TargetAppID", in.GetTargetAppID(), "error", err)
		return &npool.GetNAppAccountsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	limit := int32(constant1.DefaultLimit)
	if in.GetLimit() > 0 {
		limit = in.GetLimit()
	}

	infos, total, err := user1.GetAppAccounts(ctx, in.GetTargetAppID(), in.GetOffset(), limit)
	if err != nil {
		logger.Sugar().Errorw("GetNAppAccounts", "error", err)
		return &npool.GetNAppAccountsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetNAppAccountsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
