package goodbenefit

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/goodbenefit"

	constant "github.com/NpoolPlatform/account-gateway/pkg/message/const"

	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	return &npool.CreateAccountResponse{}, status.Error(codes.Unimplemented, "not implemented")
}
