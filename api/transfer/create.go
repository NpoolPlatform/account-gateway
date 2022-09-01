package transfer

import (
	"context"

	constant "github.com/NpoolPlatform/account-gateway/pkg/message/const"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool/account/gw/v1/transfer"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	mtransfer "github.com/NpoolPlatform/account-gateway/pkg/transfer"
)

func (s *Server) CreateTransfer(
	ctx context.Context,
	in *transfer.CreateTransferRequest,
) (
	resp *transfer.CreateTransferResponse,
	err error,
) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateTransfer")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	err = validate(in)
	if err != nil {
		return nil, err
	}

	info, err := mtransfer.CreateTransfer(
		ctx,
		in.GetAppID(),
		in.GetUserID(),
		in.GetAccount(),
		in.GetAccountType(),
		in.VerificationCode,
		in.GetTargetAccount(),
		in.GetTargetAccountType(),
	)
	if err != nil {
		logger.Sugar().Errorw("CreateTransfer", "error", err)
		return &transfer.CreateTransferResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &transfer.CreateTransferResponse{
		Info: info,
	}, nil
}
