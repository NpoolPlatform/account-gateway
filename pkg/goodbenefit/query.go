package goodbenefit

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/goodbenefit"

	gbmwcli "github.com/NpoolPlatform/account-middleware/pkg/client/goodbenefit"

	goodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/good"

	coininfocli "github.com/NpoolPlatform/chain-middleware/pkg/client/coin"
)

func GetAccount(ctx context.Context, id string) (*npool.Account, error) {
	info, err := gbmwcli.GetAccount(ctx, id)
	if err != nil {
		return nil, err
	}

	good, err := goodmwcli.GetGood(ctx, info.GoodID)
	if err != nil {
		return nil, err
	}

	coin, err := coininfocli.GetCoin(ctx, info.CoinTypeID)
	if err != nil {
		return nil, err
	}

	account := &npool.Account{
		ID:         info.ID,
		GoodID:     info.GoodID,
		GoodName:   good.Title,
		GoodUnit:   good.Unit,
		CoinTypeID: info.CoinTypeID,
		CoinName:   coin.Name,
		CoinUnit:   coin.Unit,
		CoinEnv:    coin.ENV,
		CoinLogo:   coin.Logo,
		AccountID:  info.AccountID,
		Backup:     info.Backup,
		Active:     info.Active,
		Locked:     info.Locked,
		LockedBy:   info.LockedBy,
		Blocked:    info.Blocked,
		CreatedAt:  info.CreatedAt,
		UpdatedAt:  info.UpdatedAt,
	}

	return account, nil
}

func GetAccounts(ctx context.Context, offset, limit int32) ([]*npool.Account, uint32, error) {
	return nil, 0, nil
}
