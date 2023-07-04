package payment

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/payment"

	payment "github.com/NpoolPlatform/account-gateway/pkg/payment"
	paymentmwcli "github.com/NpoolPlatform/account-middleware/pkg/client/payment"

	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
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

	account, err := paymentmwcli.GetAccount(ctx, in.GetID())
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
	}
	if in.GetActive() {
		in.Blocked = &flag
	}

	info, err := payment.UpdateAccount(ctx, in.GetID(), in.Active, in.Blocked, in.Locked)
	if err != nil {
		logger.Sugar().Errorw("UpdateAccount", "error", err)
		return &npool.UpdateAccountResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateAccountResponse{
		Info: info,
	}, nil
}
