package goodbenefit

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/goodbenefit"

	gb "github.com/NpoolPlatform/account-gateway/pkg/goodbenefit"

	gbmwcli "github.com/NpoolPlatform/account-middleware/pkg/client/goodbenefit"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

func (s *Server) UpdateAccount(ctx context.Context, in *npool.UpdateAccountRequest) (*npool.UpdateAccountResponse, error) {
	var err error

	if _, err := uuid.Parse(in.GetID()); err != nil {
		logger.Sugar().Errorw("UpdateAccount", "ID", in.GetID(), "error", err)
		return &npool.UpdateAccountResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if in.GetLocked() {
		logger.Sugar().Errorw("UpdateAccount", "Locked", in.GetLocked(), "error", "cannot lock account")
		return &npool.UpdateAccountResponse{}, status.Error(codes.InvalidArgument, "cannot lock account")
	}

	account, err := gbmwcli.GetAccount(ctx, in.GetID())
	if err != nil {
		return nil, err
	}
	if account.Blocked && (in.Blocked == nil || in.GetBlocked()) {
		logger.Sugar().Errorw("UpdateAccount", "Blocked", in.GetBlocked(), "error", "can not make change when account is blocked")
		return &npool.UpdateAccountResponse{}, status.Error(codes.InvalidArgument, "can not make change when account is blocked")
	}

	flag := false
	if in.GetBlocked() {
		in.Active = &flag

		trueFlag := true // blocked为true,backup为true,active为false
		in.Backup = &trueFlag
	}
	if in.GetActive() {
		in.Blocked = &flag
	}

	info, err := gb.UpdateAccount(ctx, in.GetID(), in.Backup, in.Active, in.Blocked, in.Locked)
	if err != nil {
		logger.Sugar().Errorw("UpdateAccount", "error", err)
		return &npool.UpdateAccountResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateAccountResponse{
		Info: info,
	}, nil
}
