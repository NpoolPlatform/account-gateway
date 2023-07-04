package user

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	commonpb "github.com/NpoolPlatform/message/npool"
	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/user"

	useraccmgrcli "github.com/NpoolPlatform/account-manager/pkg/client/user"
	useraccmgrpb "github.com/NpoolPlatform/message/npool/account/mgr/v1/user"

	user1 "github.com/NpoolPlatform/account-gateway/pkg/user"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

func (s *Server) DeleteAccount(
	ctx context.Context,
	in *npool.DeleteAccountRequest,
) (
	*npool.DeleteAccountResponse,
	error,
) {
	var err error

	if _, err := uuid.Parse(in.GetID()); err != nil {
		logger.Sugar().Errorw("DeleteAccount", "ID", in.GetID(), "error", err)
		return &npool.DeleteAccountResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("DeleteAccount", "AppID", in.GetAppID(), "error", err)
		return &npool.DeleteAccountResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if _, err := uuid.Parse(in.GetUserID()); err != nil {
		logger.Sugar().Errorw("DeleteAccount", "UserID", in.GetUserID(), "error", err)
		return &npool.DeleteAccountResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	exist, err := useraccmgrcli.ExistAccountConds(ctx, &useraccmgrpb.Conds{
		ID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetID(),
		},
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetAppID(),
		},
		UserID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetUserID(),
		},
	})
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteAccount",
			"ID", in.GetID(),
			"AppID", in.GetAppID(),
			"UserID", in.GetUserID(),
			"error", err,
		)
		return &npool.DeleteAccountResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if !exist {
		logger.Sugar().Errorw(
			"DeleteAccount",
			"ID", in.GetID(),
			"AppID", in.GetAppID(),
			"UserID", in.GetUserID(),
			"error", err,
		)
		return &npool.DeleteAccountResponse{}, status.Error(codes.InvalidArgument, "invalid account")
	}

	info, err := user1.DeleteAccount(ctx, in.GetID())
	if err != nil {
		logger.Sugar().Errorw("DeleteAccount", "error", err)
		return &npool.DeleteAccountResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.DeleteAccountResponse{
		Info: info,
	}, nil
}
