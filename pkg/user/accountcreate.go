package user

import (
	"context"
	"fmt"
	"strings"

	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/user"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	coininfocli "github.com/NpoolPlatform/chain-middleware/pkg/client/coin"

	sphinxproxypb "github.com/NpoolPlatform/message/npool/sphinxproxy"
	sphinxproxycli "github.com/NpoolPlatform/sphinx-proxy/pkg/client"

	useraccmwcli "github.com/NpoolPlatform/account-middleware/pkg/client/user"
	useraccmwpb "github.com/NpoolPlatform/message/npool/account/mw/v1/user"

	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	usercodemwcli "github.com/NpoolPlatform/basal-middleware/pkg/client/usercode"
	usercodemwpb "github.com/NpoolPlatform/message/npool/basal/mw/v1/usercode"
)

type createHandler struct {
	*Handler
	coinName *string
}

func (h *createHandler) validate(ctx context.Context) error { //nolint
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
	if h.CoinTypeID == nil {
		return fmt.Errorf("invalid coinTypeID")
	}
	if h.Address == nil || *h.Address == "" {
		return fmt.Errorf("invalid address")
	}
	if h.UsedFor == nil {
		return fmt.Errorf("invalid usedFor")
	}

	switch *h.UsedFor {
	case basetypes.AccountUsedFor_UserWithdraw:
	case basetypes.AccountUsedFor_UserDirectBenefit:
	default:
		return fmt.Errorf("usedFor %v invalid", *h.UsedFor)
	}

	return nil
}

func (h *createHandler) checkVerifyUserCode(ctx context.Context) error {
	user, err := usermwcli.GetUser(ctx, *h.AppID, *h.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("invalid user")
	}

	if *h.AccountType == basetypes.SignMethod_Google {
		h.Account = &user.GoogleSecret
	}

	if err := usercodemwcli.VerifyUserCode(ctx, &usercodemwpb.VerifyUserCodeRequest{
		Prefix:      basetypes.Prefix_PrefixUserCode.String(),
		AppID:       *h.AppID,
		Account:     *h.Account,
		AccountType: *h.AccountType,
		UsedFor:     basetypes.UsedFor_SetWithdrawAddress,
		Code:        *h.VerificationCode,
	}); err != nil {
		return err
	}

	return nil
}

func (h *createHandler) getCoinName(ctx context.Context) error {
	coin, err := coininfocli.GetCoin(ctx, *h.CoinTypeID)
	if err != nil {
		return err
	}
	if coin == nil {
		return fmt.Errorf("invlaid coin")
	}
	h.coinName = &coin.Name

	return nil
}

func (h *createHandler) checkAddress(ctx context.Context) error {
	if !strings.Contains(*h.coinName, "ironfish") {
		bal, err := sphinxproxycli.GetBalance(ctx, &sphinxproxypb.GetBalanceRequest{
			Name:    *h.coinName,
			Address: *h.Address,
		})
		if err != nil {
			return err
		}
		if bal == nil {
			return fmt.Errorf("invalid address")
		}
	}

	return nil
}

func (h *Handler) CreateAccount(ctx context.Context) (*npool.Account, error) {
	if h.AppID == nil {
		return nil, fmt.Errorf("invalid appID")
	}
	handler := &createHandler{
		Handler: h,
	}
	if err := handler.validate(ctx); err != nil {
		return nil, err
	}

	info, err := useraccmwcli.CreateAccount(ctx, &useraccmwpb.AccountReq{
		AppID:      h.AppID,
		UserID:     h.UserID,
		CoinTypeID: h.CoinTypeID,
		Address:    h.Address,
		UsedFor:    h.UsedFor,
		Labels:     *h.Labels,
		Memo:       h.Memo,
	})
	if err != nil {
		return nil, err
	}
	h.ID = &info.ID
	if err := handler.checkVerifyUserCode(ctx); err != nil {
		return nil, err
	}
	if err := handler.getCoinName(ctx); err != nil {
		return nil, err
	}
	if err := handler.checkAddress(ctx); err != nil {
		return nil, err
	}

	return h.GetAccount(ctx)
}
