//nolint:dupl
package transfer

import (
	"context"

	constant "github.com/NpoolPlatform/account-gateway/pkg/message/const"
	commontracer "github.com/NpoolPlatform/account-gateway/pkg/tracer"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool/account/gw/v1/transfer"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	mtransfer "github.com/NpoolPlatform/account-gateway/pkg/transfer"
)

func (s *Server) GetTransfers(ctx context.Context, in *transfer.GetTransfersRequest) (resp *transfer.GetTransfersResponse, err error) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetTransfers")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("GetTransfers", "AppID", in.GetAppID(), "error", err)
		return &transfer.GetTransfersResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	span = commontracer.TraceInvoker(span, "transfer", "transfer", "GetTransfers")

	infos, total, err := mtransfer.GetTransfers(ctx, in.GetAppID(), in.GetUserID(), in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorw("GetTransfers", "error", err)
		return &transfer.GetTransfersResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &transfer.GetTransfersResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetAppTransfers(
	ctx context.Context,
	in *transfer.GetAppTransfersRequest,
) (
	resp *transfer.GetAppTransfersResponse,
	err error,
) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetTransfers")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("GetTransfers", "AppID", in.GetAppID(), "error", err)
		return &transfer.GetAppTransfersResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := mtransfer.GetAppTransfers(ctx, in.GetAppID(), in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorw("GetTransfers", "error", err)
		return &transfer.GetAppTransfersResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &transfer.GetAppTransfersResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetNAppTransfers(
	ctx context.Context,
	in *transfer.GetNAppTransfersRequest,
) (
	resp *transfer.GetNAppTransfersResponse,
	err error,
) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetTransfers")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetTargetAppID()); err != nil {
		logger.Sugar().Errorw("GetNAppTransfers", "TargetAppID", in.GetTargetAppID(), "error", err)
		return &transfer.GetNAppTransfersResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := mtransfer.GetAppTransfers(ctx, in.GetTargetAppID(), in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorw("GetNAppTransfers", "error", err)
		return &transfer.GetNAppTransfersResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &transfer.GetNAppTransfersResponse{
		Infos: infos,
		Total: total,
	}, nil
}
