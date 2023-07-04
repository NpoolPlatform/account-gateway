package transfer

import (
	"context"

	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool/account/gw/v1/transfer"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	mtransfer "github.com/NpoolPlatform/account-gateway/pkg/transfer"

	"github.com/google/uuid"
)

func (s *Server) CreateTransfer(
	ctx context.Context,
	in *transfer.CreateTransferRequest,
) (
	resp *transfer.CreateTransferResponse,
	err error,
) {
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("CreateTransfer", "AppID", in.GetAppID(), "error", err)
		return &transfer.CreateTransferResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if _, err := uuid.Parse(in.GetUserID()); err != nil {
		logger.Sugar().Errorw("CreateTransfer", "UserID", in.GetUserID(), "error", err)
		return &transfer.CreateTransferResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	switch in.GetAccountType() {
	case basetypes.SignMethod_Email, basetypes.SignMethod_Mobile:
		if in.GetAccount() == "" {
			logger.Sugar().Errorw("CreateTransfer", "Account empty", "Account", in.GetAccount())
			return &transfer.CreateTransferResponse{}, status.Error(codes.InvalidArgument, "Account id empty")
		}
	case basetypes.SignMethod_Google:
	default:
		logger.Sugar().Errorw("CreateTransfer", "AccountType empty", "AccountType", in.GetAccountType())
		return &transfer.CreateTransferResponse{}, status.Error(codes.InvalidArgument, "AccountType id invalid")
	}

	if in.GetVerificationCode() == "" {
		logger.Sugar().Errorw("CreateTransfer", "VerificationCode empty", "VerificationCode", in.GetVerificationCode())
		return &transfer.CreateTransferResponse{}, status.Error(codes.InvalidArgument, "VerificationCode id empty")
	}
	if in.GetTargetAccount() == "" {
		logger.Sugar().Errorw("CreateTransfer", "TargetAccount empty", "TargetAccount", in.GetTargetAccount())
		return &transfer.CreateTransferResponse{}, status.Error(codes.InvalidArgument, "TargetAccount id empty")
	}

	switch in.GetTargetAccountType() {
	case basetypes.SignMethod_Email:
	case basetypes.SignMethod_Mobile:
	default:
		logger.Sugar().Errorw("CreateTransfer", "TargetAccountType empty", "TargetAccountType", in.GetTargetAccountType())
		return &transfer.CreateTransferResponse{}, status.Error(codes.InvalidArgument, "TargetAccountType id invalid")
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
