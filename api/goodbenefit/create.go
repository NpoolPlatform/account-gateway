package goodbenefit

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/goodbenefit"

	constant "github.com/NpoolPlatform/account-gateway/pkg/message/const"

	gb "github.com/NpoolPlatform/account-gateway/pkg/goodbenefit"

	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

func (s *Server) CreateAccount(ctx context.Context, in *npool.CreateAccountRequest) (*npool.CreateAccountResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateAccountAccount")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetGoodID()); err != nil {
		return &npool.CreateAccountResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := gb.CreateAccount(ctx, in.GetGoodID())
	if err != nil {
		return &npool.CreateAccountResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateAccountResponse{
		Info: info,
	}, nil
}
