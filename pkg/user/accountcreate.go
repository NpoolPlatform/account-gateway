package user

import (
	"context"
	"fmt"
	"strings"

	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/user"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	coininfocli "github.com/NpoolPlatform/chain-middleware/pkg/client/coin"

	sphinxproxypb "github.com/NpoolPlatform/message/npool/sphinxproxy"
	sphinxproxycli "github.com/NpoolPlatform/sphinx-proxy/pkg/client"

	useraccmwcli "github.com/NpoolPlatform/account-middleware/pkg/client/user"
	useraccmwpb "github.com/NpoolPlatform/message/npool/account/mw/v1/user"

	accountmgrpb "github.com/NpoolPlatform/message/npool/account/mgr/v1/account"

	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	usercodemwcli "github.com/NpoolPlatform/basal-middleware/pkg/client/usercode"
	usercodemwpb "github.com/NpoolPlatform/message/npool/basal/mw/v1/usercode"
)

func CreateAccount(
	ctx context.Context,
	appID, userID, coinTypeID string,
	usedFor accountmgrpb.AccountUsedFor,
	address string,
	labels []string,
	account string,
	accountType basetypes.SignMethod,
	verificationCode string,
) (
	*npool.Account, error,
) {
	a, err := appmwcli.GetApp(ctx, appID)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, fmt.Errorf("invalid app")
	}

	u, err := usermwcli.GetUser(ctx, appID, userID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, fmt.Errorf("invalid user")
	}

	if accountType == basetypes.SignMethod_Google {
		account = u.GoogleSecret
	}

	if err := usercodemwcli.VerifyUserCode(ctx, &usercodemwpb.VerifyUserCodeRequest{
		Prefix:      basetypes.Prefix_PrefixUserCode.String(),
		AppID:       appID,
		Account:     account,
		AccountType: accountType,
		UsedFor:     basetypes.UsedFor_SetWithdrawAddress,
		Code:        verificationCode,
	}); err != nil {
		return nil, err
	}

	coin, err := coininfocli.GetCoin(ctx, coinTypeID)
	if err != nil {
		return nil, err
	}
	if coin == nil {
		return nil, fmt.Errorf("invlaid coin")
	}

	if !strings.Contains(coin.Name, "ironfish") {
		bal, err := sphinxproxycli.GetBalance(ctx, &sphinxproxypb.GetBalanceRequest{
			Name:    coin.Name,
			Address: address,
		})
		if err != nil {
			return nil, err
		}
		if bal == nil {
			return nil, fmt.Errorf("invalid address")
		}
	}

	info, err := useraccmwcli.CreateAccount(ctx, &useraccmwpb.AccountReq{
		AppID:      &appID,
		UserID:     &userID,
		CoinTypeID: &coinTypeID,
		Address:    &address,
		UsedFor:    &usedFor,
		Labels:     labels,
	})
	if err != nil {
		return nil, err
	}

	return GetAccount(ctx, info.ID)
}
