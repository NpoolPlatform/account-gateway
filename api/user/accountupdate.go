package user

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/user"

	constant "github.com/NpoolPlatform/account-gateway/pkg/message/const"
	useraccmwcli "github.com/NpoolPlatform/account-middleware/pkg/client/user"

	user1 "github.com/NpoolPlatform/account-gateway/pkg/user"

	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

func (s *Server) UpdateAccount(
	ctx context.Context,
	in *npool.UpdateAccountRequest,
) (
	*npool.UpdateAccountResponse,
	error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "UpdateAccount")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetID()); err != nil {
		logger.Sugar().Errorw("UpdateAccount", "ID", in.GetID(), "error", err)
		return &npool.UpdateAccountResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("UpdateAccount", "AppID", in.GetAppID(), "error", err)
		return &npool.UpdateAccountResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if _, err := uuid.Parse(in.GetUserID()); err != nil {
		logger.Sugar().Errorw("UpdateAccount", "UserID", in.GetUserID(), "error", err)
		return &npool.UpdateAccountResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	account, err := useraccmwcli.GetAccount(ctx, in.GetID())
	if err != nil {
		logger.Sugar().Errorw("UpdateAccount", "ID", in.GetID(), "error", err)
		return nil, err
	}
	if account.AppID != in.GetAppID() {
		logger.Sugar().Errorw("UpdateAccount", "AppID", in.GetAppID(), "error", "Wrong AppID")
		return &npool.UpdateAccountResponse{}, status.Error(codes.InvalidArgument, "Wrong AppID")
	}
	if account.UserID != in.GetUserID() {
		logger.Sugar().Errorw("UpdateAccount", "UserID", in.GetUserID(), "error", "Wrong UserID")
		return &npool.UpdateAccountResponse{}, status.Error(codes.InvalidArgument, "Wrong UserID")
	}

	info, err := user1.UpdateAccount(ctx, in.GetID(), nil, nil, in.Labels)
	if err != nil {
		logger.Sugar().Errorw("UpdateAccount", "error", err)
		return &npool.UpdateAccountResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateAccountResponse{
		Info: info,
	}, nil
}

func (s *Server) UpdateAppUserAccount(
	ctx context.Context,
	in *npool.UpdateAppUserAccountRequest,
) (
	*npool.UpdateAppUserAccountResponse,
	error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "UpdateAppUserAccount")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetID()); err != nil {
		logger.Sugar().Errorw("UpdateAppUserAccount", "ID", in.GetID(), "error", err)
		return &npool.UpdateAppUserAccountResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if _, err := uuid.Parse(in.GetTargetAppID()); err != nil {
		logger.Sugar().Errorw("UpdateAppUserAccount", "TargetAppID", in.GetTargetAppID(), "error", err)
		return &npool.UpdateAppUserAccountResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if _, err := uuid.Parse(in.GetTargetUserID()); err != nil {
		logger.Sugar().Errorw("UpdateAppUserAccount", "TargetUserID", in.GetTargetUserID(), "error", err)
		return &npool.UpdateAppUserAccountResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	account, err := useraccmwcli.GetAccount(ctx, in.GetID())
	if err != nil {
		logger.Sugar().Errorw("UpdateAccount", "ID", in.GetID(), "error", err)
		return nil, err
	}
	if account.AppID != in.GetTargetAppID() {
		logger.Sugar().Errorw("UpdateAppUserAccount", "TargetAppID", in.GetTargetAppID(), "error", "Wrong TargetAppID")
		return &npool.UpdateAppUserAccountResponse{}, status.Error(codes.InvalidArgument, "Wrong TargetAppID")
	}
	if account.UserID != in.GetTargetUserID() {
		logger.Sugar().Errorw("UpdateAppUserAccount", "TargetUserID", in.GetTargetUserID(), "error", "Wrong TargetUserID")
		return &npool.UpdateAppUserAccountResponse{}, status.Error(codes.InvalidArgument, "Wrong TargetUserID")
	}

	if account.Blocked && in.GetActive() {
		logger.Sugar().Errorw("UpdateAppUserAccount", "Active", in.GetActive(), "error", "Account is blocked")
		return &npool.UpdateAppUserAccountResponse{}, status.Error(codes.InvalidArgument, "Account is blocked")
	}

	if account.Blocked && (in.Blocked == nil || in.GetBlocked()) {
		info, err := user1.GetAccount(ctx, in.GetID())
		if err != nil {
			logger.Sugar().Errorw("UpdateAppUserAccount", "error", err)
			return &npool.UpdateAppUserAccountResponse{}, status.Error(codes.InvalidArgument, err.Error())
		}
		return &npool.UpdateAppUserAccountResponse{
			Info: info,
		}, nil
	}

	falseFlag := false
	if in.GetBlocked() {
		in.Active = &falseFlag
	}

	info, err := user1.UpdateAccount(ctx, in.GetID(), in.Active, in.Blocked, nil)
	if err != nil {
		logger.Sugar().Errorw("UpdateAppUserAccount", "error", err)
		return &npool.UpdateAppUserAccountResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateAppUserAccountResponse{
		Info: info,
	}, nil
}
