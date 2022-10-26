package user

import (
	"context"
	"fmt"

	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/user"

	useraccmwcli "github.com/NpoolPlatform/account-middleware/pkg/client/user"
	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	coininfocli "github.com/NpoolPlatform/sphinx-coininfo/pkg/client"
)

func GetWithdrawAccount(ctx context.Context, id string) (*npool.Account, error) {
	info, err := useraccmwcli.GetAccount(ctx, id)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf("invalid account")
	}

	u, err := usermwcli.GetUser(ctx, info.AppID, info.UserID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, fmt.Errorf("invalid user")
	}

	coin, err := coininfocli.GetCoinInfo(ctx, info.CoinTypeID)
	if err != nil {
		return nil, err
	}
	if coin == nil {
		return nil, fmt.Errorf("invalid coin")
	}

	acc := &npool.Account{
		ID:           info.ID,
		AppID:        info.AppID,
		UserID:       info.UserID,
		CoinTypeID:   info.CoinTypeID,
		CoinName:     coin.Name,
		CoinUnit:     coin.Unit,
		CoinEnv:      coin.ENV,
		CoinLogo:     coin.Logo,
		AccountID:    info.AccountID,
		Address:      info.Address,
		UsedFor:      info.UsedFor,
		CreatedAt:    info.CreatedAt,
		PhoneNO:      u.PhoneNO,
		EmailAddress: u.EmailAddress,
	}

	return acc, nil
}
