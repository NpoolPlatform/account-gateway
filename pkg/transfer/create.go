package transfer

import (
	"context"
	"fmt"

	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	usercodemwcli "github.com/NpoolPlatform/basal-middleware/pkg/client/usercode"
	usercodemwpb "github.com/NpoolPlatform/message/npool/basal/mw/v1/usercode"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/transfer"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"

	transfermwcli "github.com/NpoolPlatform/account-middleware/pkg/client/transfer"
	transfermwpb "github.com/NpoolPlatform/message/npool/account/mw/v1/transfer"
)

type createHandler struct {
	*Handler
	infos []*transfermwpb.Transfer
	users map[string]*usermwpb.User
	accs  []*npool.Transfer
}

func (h *createHandler) validate() error {
	if h.AppID == nil {
		return fmt.Errorf("invalid appID")
	}
	if h.UserID == nil {
		return fmt.Errorf("invalid userID")
	}
	if h.AccountType == nil {
		return fmt.Errorf("invalid accountType")
	}
	switch *h.AccountType {
	case basetypes.SignMethod_Email:
		fallthrough //nolint
	case basetypes.SignMethod_Mobile:
		if h.Account == nil || *h.Account == "" {
			return fmt.Errorf("account is empty")
		}
	case basetypes.SignMethod_Google:
	default:
		return fmt.Errorf("accountType %v invalid", *h.AccountType)
	}
	if h.VerificationCode == nil || *h.VerificationCode == "" {
		return fmt.Errorf("invalid verificationCode")
	}

	if h.TargetAccountType == nil {
		return fmt.Errorf("invalid targetAccountType")
	}
	switch *h.TargetAccountType {
	case basetypes.SignMethod_Email:
	case basetypes.SignMethod_Mobile:
	default:
		return fmt.Errorf("targetAccountType %v invalid", *h.TargetAccountType)
	}
	return nil
}

func (h *createHandler) formalize() {
	for _, val := range h.infos {
		userInfo, ok := h.users[val.TargetUserID]
		if !ok {
			continue
		}

		h.accs = append(h.accs, &npool.Transfer{
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
}

func (h *createHandler) getUsers(ctx context.Context) error {
	userInfo, err := usermwcli.GetUser(ctx, *h.AppID, *h.UserID)
	if err != nil {
		return err
	}

	if userInfo == nil {
		return fmt.Errorf("user not found")
	}

	if *h.AccountType == basetypes.SignMethod_Google {
		h.Account = &userInfo.GoogleSecret
	}

	if err := usercodemwcli.VerifyUserCode(ctx, &usercodemwpb.VerifyUserCodeRequest{
		Prefix:      basetypes.Prefix_PrefixUserCode.String(),
		AppID:       *h.AppID,
		Account:     *h.Account,
		AccountType: *h.AccountType,
		UsedFor:     basetypes.UsedFor_SetTransferTargetUser,
		Code:        *h.VerificationCode,
	}); err != nil {
		return err
	}

	conds := &usermwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	}
	switch *h.TargetAccountType {
	case basetypes.SignMethod_Email:
		conds.EmailAddress = &basetypes.StringVal{Op: cruder.EQ, Value: *h.TargetAccount}
	case basetypes.SignMethod_Mobile:
		conds.PhoneNO = &basetypes.StringVal{Op: cruder.EQ, Value: *h.TargetAccount}
	}

	targetUser, err := usermwcli.GetUserOnly(ctx, conds)
	if err != nil {
		return err
	}
	if targetUser == nil {
		return fmt.Errorf("target user not found")
	}

	if *h.UserID == targetUser.ID {
		return fmt.Errorf("cannot set yourself as the payee")
	}
	h.users[targetUser.ID] = targetUser
	h.TargetUserID = &targetUser.ID

	return nil
}

func (h *Handler) CreateTransfer(ctx context.Context) (*npool.Transfer, error) {
	handler := &createHandler{
		Handler: h,
		infos:   []*transfermwpb.Transfer{},
		users:   map[string]*usermwpb.User{},
	}
	if err := handler.validate(); err != nil {
		return nil, err
	}
	if err := handler.getUsers(ctx); err != nil {
		return nil, err
	}

	exist, err := transfermwcli.ExistTransferConds(ctx, &transfermwpb.Conds{
		AppID: &basetypes.StringVal{
			Op:    cruder.EQ,
			Value: *h.AppID,
		},
		UserID: &basetypes.StringVal{
			Op:    cruder.EQ,
			Value: *h.UserID,
		},
		TargetUserID: &basetypes.StringVal{
			Op:    cruder.EQ,
			Value: *handler.TargetUserID,
		},
	})
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, fmt.Errorf("target user already exist")
	}

	info, err := transfermwcli.CreateTransfer(ctx, &transfermwpb.TransferReq{
		AppID:        h.AppID,
		UserID:       h.UserID,
		TargetUserID: h.TargetUserID,
	})
	if err != nil {
		return nil, err
	}
	handler.infos = append(handler.infos, info)
	handler.formalize()

	return handler.accs[0], nil
}
