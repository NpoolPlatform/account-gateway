package platform

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/platform"
	accountmgrpb "github.com/NpoolPlatform/message/npool/account/mgr/v1/account"

	platform1 "github.com/NpoolPlatform/account-gateway/pkg/platform"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

func (s *Server) CreateAccount(ctx context.Context, in *npool.CreateAccountRequest) (*npool.CreateAccountResponse, error) {
	var err error

	if _, err := uuid.Parse(in.GetCoinTypeID()); err != nil {
		logger.Sugar().Errorw("CreateAccount", "CoinTypeID", in.GetCoinTypeID(), "error", err)
		return &npool.CreateAccountResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if in.Address != nil && in.GetAddress() == "" {
		logger.Sugar().Errorw("CreateAccount", "Address", in.GetAddress())
		return &npool.CreateAccountResponse{}, status.Error(codes.InvalidArgument, "invalid address")
	}
	switch in.GetUsedFor() {
	case accountmgrpb.AccountUsedFor_UserBenefitHot:
		fallthrough // nolint
	case accountmgrpb.AccountUsedFor_GasProvider:
		if in.Address != nil {
			logger.Sugar().Errorw("CreateAccount", "Address", in.GetAddress())
			return &npool.CreateAccountResponse{}, status.Error(codes.InvalidArgument, "invalid account")
		}
	case accountmgrpb.AccountUsedFor_PaymentCollector:
		fallthrough // nolint
	case accountmgrpb.AccountUsedFor_UserBenefitCold:
		fallthrough // nolint
	case accountmgrpb.AccountUsedFor_PlatformBenefitCold:
		if in.Address == nil {
			logger.Sugar().Errorw("CreateAccount", "Address", in.GetAddress())
			return &npool.CreateAccountResponse{}, status.Error(codes.InvalidArgument, "invalid account")
		}
	}

	info, err := platform1.CreateAccount(ctx, in.GetCoinTypeID(), in.Address, in.GetUsedFor())
	if err != nil {
		logger.Sugar().Errorw("CreateAccount", "error", err)
		return &npool.CreateAccountResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateAccountResponse{
		Info: info,
	}, nil
}
