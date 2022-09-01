package transfer

import (
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool/account/gw/v1/transfer"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	signmethodpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/signmethod"
)

func validate(in *transfer.CreateTransferRequest) (err error) {
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("GetDepositAccount", "AppID", in.GetAppID(), "error", err)
		return status.Error(codes.InvalidArgument, err.Error())
	}
	if _, err := uuid.Parse(in.GetUserID()); err != nil {
		logger.Sugar().Errorw("GetDepositAccount", "UserID", in.GetUserID(), "error", err)
		return status.Error(codes.InvalidArgument, err.Error())
	}
	if in.GetAccount() == "" {
		logger.Sugar().Errorw("GetDepositAccount", "Account empty", "Account", in.GetAccount())
		return status.Error(codes.InvalidArgument, "Account id empty")
	}

	switch in.GetAccountType() | in.GetTargetAccountType() {
	case signmethodpb.SignMethodType_Email:
	case signmethodpb.SignMethodType_Mobile:
	default:
		logger.Sugar().Errorw("GetDepositAccount", "AccountType empty", "AccountType", in.GetAccountType())
		return status.Error(codes.InvalidArgument, "AccountType id invalid")
	}

	if in.GetVerificationCode() == "" {
		logger.Sugar().Errorw("GetDepositAccount", "VerificationCode empty", "VerificationCode", in.GetVerificationCode())
		return status.Error(codes.InvalidArgument, "VerificationCode id empty")
	}
	if in.GetTargetAccount() == "" {
		logger.Sugar().Errorw("GetDepositAccount", "TargetAccount empty", "TargetAccount", in.GetTargetAccount())
		return status.Error(codes.InvalidArgument, "TargetAccount id empty")
	}

	return nil
}
