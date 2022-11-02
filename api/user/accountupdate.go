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

	// TODO: query id belong to AppID / UserID

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

	info, err := user1.UpdateAccount(ctx, in.GetID(), in.Active, in.Blocked, nil)
	if err != nil {
		logger.Sugar().Errorw("UpdateAppUserAccount", "error", err)
		return &npool.UpdateAppUserAccountResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateAppUserAccountResponse{
		Info: info,
	}, nil
}