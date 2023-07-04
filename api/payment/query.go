package payment

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	constant1 "github.com/NpoolPlatform/account-gateway/pkg/const"

	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/payment"

	payment1 "github.com/NpoolPlatform/account-gateway/pkg/payment"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetAccounts(ctx context.Context, in *npool.GetAccountsRequest) (*npool.GetAccountsResponse, error) {
	var err error

	limit := int32(constant1.DefaultLimit)
	if in.GetLimit() > 0 {
		limit = in.GetLimit()
	}

	infos, total, err := payment1.GetAccounts(ctx, in.GetOffset(), limit)
	if err != nil {
		logger.Sugar().Errorw("GetAccounts", "Offset", in.GetOffset(), "Limit", limit, "error", err)
		return &npool.GetAccountsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAccountsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
