package transfer

import (
	"context"
	"fmt"

	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	usercodemwcli "github.com/NpoolPlatform/basal-middleware/pkg/client/usercode"
	usercodemwpb "github.com/NpoolPlatform/message/npool/basal/mw/v1/usercode"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"

	"github.com/NpoolPlatform/message/npool/account/gw/v1/transfer"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"

	mgrcli "github.com/NpoolPlatform/account-manager/pkg/client/transfer"
	mgrpb "github.com/NpoolPlatform/message/npool/account/mgr/v1/transfer"
)

// nolint:funlen
func CreateTransfer(ctx context.Context,
	appID,
	userID,
	account string,
	accountType basetypes.SignMethod,
	verificationCode,
	targetAccount string,
	targetAccountType basetypes.SignMethod,
) (
	*transfer.Transfer, error,
) {
	var err error

	userInfo, err := usermwcli.GetUser(ctx, appID, userID)
	if err != nil {
		return nil, err
	}

	if userInfo == nil {
		return nil, fmt.Errorf("user not found")
	}

	if accountType == basetypes.SignMethod_Google {
		account = userInfo.GoogleSecret
	}

	if err := usercodemwcli.VerifyUserCode(ctx, &usercodemwpb.VerifyUserCodeRequest{
		Prefix:      basetypes.Prefix_PrefixUserCode.String(),
		AppID:       appID,
		Account:     account,
		AccountType: accountType,
		UsedFor:     basetypes.UsedFor_SetTransferTargetUser,
		Code:        verificationCode,
	}); err != nil {
		return nil, err
	}

	conds := &usermwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: appID},
	}
	switch targetAccountType {
	case basetypes.SignMethod_Email:
		conds.EmailAddress = &basetypes.StringVal{Op: cruder.EQ, Value: targetAccount}
	case basetypes.SignMethod_Mobile:
		conds.PhoneNO = &basetypes.StringVal{Op: cruder.EQ, Value: targetAccount}
	}

	targetUser, err := usermwcli.GetUserOnly(ctx, conds)
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

	info, err := mgrcli.CreateTransfer(ctx, &mgrpb.TransferReq{
		AppID:        &appID,
		UserID:       &userID,
		TargetUserID: &targetUser.ID,
	})
	if err != nil {
		return nil, err
	}

	targetUserInfo, err := usermwcli.GetUser(ctx, appID, targetUser.ID)
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

	targetUser, err := usermwcli.GetUser(ctx, info.AppID, info.TargetUserID)
	if err != nil {
		return nil, err
	}

	if targetUser == nil {
		targetUser = &usermwpb.User{}
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

	for _, val := range infos {
		targetUserIDs = append(targetUserIDs, val.TargetUserID)
	}

	users, _, err := usermwcli.GetUsers(ctx, &usermwpb.Conds{
		IDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: targetUserIDs},
	}, 0, int32(len(targetUserIDs)))
	if err != nil {
		return nil, err
	}
	targetUser := map[string]*usermwpb.User{}

	for _, val := range users {
		targetUser[val.ID] = val
	}

	transferInfos := []*transfer.Transfer{}

	for _, val := range infos {
		userInfo := &usermwpb.User{}

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
