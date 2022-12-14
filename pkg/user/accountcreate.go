package user

import (
	"context"

	"fmt"

	"github.com/NpoolPlatform/message/npool/appuser/mgr/v2/signmethod"

	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/user"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	coininfocli "github.com/NpoolPlatform/chain-middleware/pkg/client/coin"

	sphinxproxypb "github.com/NpoolPlatform/message/npool/sphinxproxy"
	sphinxproxycli "github.com/NpoolPlatform/sphinx-proxy/pkg/client"

	useraccmwcli "github.com/NpoolPlatform/account-middleware/pkg/client/user"
	useraccmwpb "github.com/NpoolPlatform/message/npool/account/mw/v1/user"

	accountmgrpb "github.com/NpoolPlatform/message/npool/account/mgr/v1/account"

	"github.com/NpoolPlatform/message/npool/third/mgr/v1/usedfor"
	thirdmwcli "github.com/NpoolPlatform/third-middleware/pkg/client/verify"
)

func CreateAccount(
	ctx context.Context,
	appID, userID, coinTypeID string,
	usedFor accountmgrpb.AccountUsedFor,
	address string,
	labels []string,
	account string,
	accountType signmethod.SignMethodType,
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

	if accountType == signmethod.SignMethodType_Google {
		account = u.GoogleSecret
	}

	if err := thirdmwcli.VerifyCode(
		ctx,
		appID,
		account,
		verificationCode,
		accountType,
		usedfor.UsedFor_SetWithdrawAddress,
	); err != nil {
		return nil, err
	}

	coin, err := coininfocli.GetCoin(ctx, coinTypeID)
	if err != nil {
		return nil, err
	}
	if coin == nil {
		return nil, fmt.Errorf("invlaid coin")
	}

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
