package transfer

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/message/npool/third/mgr/v1/usedfor"
	thirdmwcli "github.com/NpoolPlatform/third-middleware/pkg/client/verify"

	constant "github.com/NpoolPlatform/account-gateway/pkg/message/const"
	commontracer "github.com/NpoolPlatform/account-gateway/pkg/tracer"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"

	"github.com/NpoolPlatform/message/npool/account/gw/v1/transfer"
	signmethodpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/signmethod"

	appusermgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/appuser"
	appusermgpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/appuser"

	appusermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	appusermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"

	mgrcli "github.com/NpoolPlatform/account-manager/pkg/client/transfer"
	mgrpb "github.com/NpoolPlatform/message/npool/account/mgr/v1/transfer"
)

// nolint:funlen
func CreateTransfer(ctx context.Context,
	appID,
	userID,
	account string,
	accountType signmethodpb.SignMethodType,
	verificationCode,
	targetAccount string,
	targetAccountType signmethodpb.SignMethodType,
) (
	*transfer.Transfer, error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateTransfer")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceInvoker(span, "transfer", "third-gateway", "VerifyCode")

	userInfo, err := appusermwcli.GetUser(ctx, appID, userID)
	if err != nil {
		return nil, err
	}

	if accountType == signmethodpb.SignMethodType_Google {
		account = userInfo.GoogleSecret
	}

	if err := thirdmwcli.VerifyCode(
		ctx,
		appID,
		account,
		verificationCode,
		accountType,
		usedfor.UsedFor_SetTransferTargetUser,
	); err != nil {
		return nil, err
	}

	conds := &appusermgpb.Conds{
		AppID: &npool.StringVal{
			Op:    cruder.EQ,
			Value: appID,
		},
	}
	switch targetAccountType {
	case signmethodpb.SignMethodType_Email:
		conds.EmailAddress = &npool.StringVal{
			Op:    cruder.EQ,
			Value: targetAccount,
		}
	case signmethodpb.SignMethodType_Mobile:
		conds.PhoneNO = &npool.StringVal{
			Op:    cruder.EQ,
			Value: targetAccount,
		}
	}

	span = commontracer.TraceInvoker(span, "transfer", "appuser-manager", "GetAppUserOnly")

	targetUser, err := appusermgrcli.GetAppUserOnly(ctx, conds)
	if err != nil {
		return nil, err
	}
	if targetUser == nil {
		return nil, fmt.Errorf("target user not found")
	}

	if userID == targetUser.ID {
		return nil, fmt.Errorf("cannot set yourself as the payee")
	}

	span = commontracer.TraceInvoker(span, "transfer", "manager", "ExistTransferConds")

	exist, err := mgrcli.ExistTransferConds(ctx, &mgrpb.Conds{
		AppID: &npool.StringVal{
			Op:    cruder.EQ,
			Value: appID,
		},
		UserID: &npool.StringVal{
			Op:    cruder.EQ,
			Value: userID,
		},
		TargetUserID: &npool.StringVal{
			Op:    cruder.EQ,
			Value: targetUser.ID,
		},
	})
	if err != nil {
		return nil, err
	}

	if exist {
		return nil, fmt.Errorf("target user already exist")
	}

	span = commontracer.TraceInvoker(span, "transfer", "manager", "CreateTransfer")

	info, err := mgrcli.CreateTransfer(ctx, &mgrpb.TransferReq{
		AppID:        &appID,
		UserID:       &userID,
		TargetUserID: &targetUser.ID,
	})
	if err != nil {
		return nil, err
	}

	targetUserInfo, err := appusermwcli.GetUser(ctx, appID, targetUser.ID)
	if err != nil {
		return nil, err
	}

	return &transfer.Transfer{
		ID:                 info.ID,
		AppID:              info.AppID,
		UserID:             info.UserID,
		TargetUserID:       info.TargetUserID,
		TargetEmailAddress: targetUserInfo.EmailAddress,
		TargetPhoneNO:      targetUserInfo.PhoneNO,
		CreatedAt:          info.CreatedAt,
		TargetUsername:     targetUserInfo.Username,
		TargetFirstName:    targetUserInfo.FirstName,
		TargetLastName:     targetUserInfo.LastName,
	}, nil
}

func DeleteTransfer(ctx context.Context, id, appID, userID string) (*transfer.Transfer, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetTransfers")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceInvoker(span, "transfer", "manager", "DeleteTransfer")

	exist, err := mgrcli.ExistTransferConds(ctx, &mgrpb.Conds{
		ID: &npool.StringVal{
			Op:    cruder.EQ,
			Value: id,
		},
		AppID: &npool.StringVal{
			Op:    cruder.EQ,
			Value: appID,
		},
		UserID: &npool.StringVal{
			Op:    cruder.EQ,
			Value: userID,
		},
	})
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("transfer not exist")
	}

	info, err := mgrcli.DeleteTransfer(ctx, id)
	if err != nil {
		return nil, err
	}

	span = commontracer.TraceInvoker(span, "transfer", "appuser-manager", "GetAppUser")

	targetUser, err := appusermwcli.GetUser(ctx, info.AppID, info.TargetUserID)
	if err != nil {
		return nil, err
	}

	if targetUser == nil {
		targetUser = &appusermwpb.User{}
	}
	return &transfer.Transfer{
		ID:                 info.ID,
		AppID:              info.AppID,
		UserID:             info.UserID,
		TargetUserID:       info.TargetUserID,
		TargetEmailAddress: targetUser.EmailAddress,
		TargetPhoneNO:      targetUser.PhoneNO,
		CreatedAt:          info.CreatedAt,
		TargetUsername:     targetUser.Username,
		TargetFirstName:    targetUser.FirstName,
		TargetLastName:     targetUser.LastName,
	}, nil
}

func GetTransfers(ctx context.Context, appID, userID string, offset, limit int32) ([]*transfer.Transfer, uint32, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetTransfers")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceInvoker(span, "transfer", "manager", "GetTransfers")

	infos, total, err := mgrcli.GetTransfers(ctx, &mgrpb.Conds{
		AppID: &npool.StringVal{
			Op:    cruder.EQ,
			Value: appID,
		},
		UserID: &npool.StringVal{
			Op:    cruder.EQ,
			Value: userID,
		},
	}, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	if len(infos) == 0 {
		return []*transfer.Transfer{}, 0, nil
	}

	transferInfos, err := expand(ctx, infos)
	if err != nil {
		return nil, 0, err
	}
	return transferInfos, total, nil
}

func GetAppTransfers(ctx context.Context, appID string, offset, limit int32) ([]*transfer.Transfer, uint32, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetTransfers")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceInvoker(span, "transfer", "manager", "GetTransfers")

	infos, total, err := mgrcli.GetTransfers(ctx, &mgrpb.Conds{
		AppID: &npool.StringVal{
			Op:    cruder.EQ,
			Value: appID,
		},
	}, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	if len(infos) == 0 {
		return []*transfer.Transfer{}, 0, nil
	}

	transferInfos, err := expand(ctx, infos)
	if err != nil {
		return nil, 0, err
	}
	return transferInfos, total, nil
}

func expand(ctx context.Context, infos []*mgrpb.Transfer) ([]*transfer.Transfer, error) {
	var err error
	targetUserIDs := []string{}

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetTransfers")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	for _, val := range infos {
		targetUserIDs = append(targetUserIDs, val.TargetUserID)
	}
	span = commontracer.TraceInvoker(span, "transfer", "appuser-manager", "GetAppUsers")

	users, _, err := appusermwcli.GetManyUsers(ctx, targetUserIDs)
	if err != nil {
		return nil, err
	}
	targetUser := map[string]*appusermwpb.User{}

	for _, val := range users {
		targetUser[val.ID] = val
	}

	transferInfos := []*transfer.Transfer{}

	for _, val := range infos {
		userInfo := &appusermwpb.User{}

		if _, ok := targetUser[val.TargetUserID]; ok {
			userInfo = targetUser[val.TargetUserID]
		}

		transferInfos = append(transferInfos, &transfer.Transfer{
			ID:                 val.ID,
			AppID:              val.AppID,
			UserID:             val.UserID,
			TargetUserID:       val.TargetUserID,
			TargetEmailAddress: userInfo.EmailAddress,
			TargetPhoneNO:      userInfo.PhoneNO,
			CreatedAt:          val.CreatedAt,
			TargetUsername:     userInfo.Username,
			TargetFirstName:    userInfo.FirstName,
			TargetLastName:     userInfo.LastName,
		})
	}
	return transferInfos, nil
}
