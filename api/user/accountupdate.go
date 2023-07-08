package user

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	user1 "github.com/NpoolPlatform/account-gateway/pkg/user"
	useraccmwcli "github.com/NpoolPlatform/account-middleware/pkg/client/user"
	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateAccount(
	ctx context.Context,
	in *npool.UpdateAccountRequest,
) (
	*npool.UpdateAccountResponse,
	error,
) {
	handler, err := user1.NewHandler(
		ctx,
		user1.WithID(&in.ID),
		user1.WithAppID(&in.AppID),
		user1.WithUserID(&in.UserID),
		user1.WithLabels(in.Labels),
		user1.WithMemo(in.Memo),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateAccount",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateAccountResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.UpdateAccount(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateAccount",
			"In", in,
			"Error", err,
		)
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
	handler, err := user1.NewHandler(
		ctx,
		user1.WithID(&in.ID),
		user1.WithAppID(&in.TargetAppID),
		user1.WithUserID(&in.TargetUserID),
		user1.WithBlocked(in.Blocked),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateAppUserAccount",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateAppUserAccountResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	account, err := useraccmwcli.GetAccount(ctx, in.GetID())
	if err != nil {
		logger.Sugar().Errorw("UpdateAccount", "ID", in.GetID(), "error", err)
		return nil, err
	}
	if account.Blocked && in.GetActive() {
		logger.Sugar().Errorw("UpdateAppUserAccount", "Active", in.GetActive(), "error", "Account is blocked")
		return &npool.UpdateAppUserAccountResponse{}, status.Error(codes.InvalidArgument, "Account is blocked")
	}

	if account.Blocked && (in.Blocked == nil || in.GetBlocked()) {
		info, err := handler.GetAccount(ctx)
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
		handler.Active = &falseFlag
	}

	info, err := handler.UpdateAccount(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateAppUserAccount",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateAppUserAccountResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateAppUserAccountResponse{
		Info: info,
	}, nil
}
