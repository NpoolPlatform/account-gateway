package contract

import (
	"context"

	contract1 "github.com/NpoolPlatform/account-gateway/pkg/contract"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/contract"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) AdminGetAccounts(ctx context.Context, in *npool.AdminGetAccountsRequest) (*npool.AdminGetAccountsResponse, error) {
	handler, err := contract1.NewHandler(
		ctx,
		contract1.WithOffset(in.GetOffset()),
		contract1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminGetAccounts",
			"In", in,
			"Error", err,
		)
		return &npool.AdminGetAccountsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetAccounts(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAccounts",
			"In", in,
			"Error", err,
		)
		return &npool.AdminGetAccountsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.AdminGetAccountsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
