package platform

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/platform"

	pltfmwcli "github.com/NpoolPlatform/account-middleware/pkg/client/platform"
	// pltfmwpb "github.com/NpoolPlatform/message/npool/account/mw/v1/platform"

	coininfocli "github.com/NpoolPlatform/sphinx-coininfo/pkg/client"
)

func GetAccount(ctx context.Context, id string) (*npool.Account, error) {
	info, err := pltfmwcli.GetAccount(ctx, id)
	if err != nil {
		return nil, err
	}

	coin, err := coininfocli.GetCoinInfo(ctx, info.CoinTypeID)
	if err != nil {
		return nil, err
	}

	account := &npool.Account{
		ID:         info.ID,
		CoinTypeID: info.CoinTypeID,
		CoinName:   coin.Name,
		CoinUnit:   coin.Unit,
		CoinEnv:    coin.ENV,
		CoinLogo:   coin.Logo,
		UsedFor:    info.UsedFor,
		AccountID:  info.AccountID,
		Address:    info.Address,
		Backup:     info.Backup,
		Active:     info.Active,
		Locked:     info.Locked,
		LockedBy:   info.LockedBy,
		Blocked:    info.Blocked,
		CreatedAt:  info.CreatedAt,
	}

	return account, nil
}

func GetAccounts(ctx context.Context, offset, limit int32) ([]*npool.Account, uint32, error) {
	return nil, 0, nil
}
