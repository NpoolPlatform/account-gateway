package payment

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/payment"

	paymentmwcli "github.com/NpoolPlatform/account-middleware/pkg/client/payment"
	paymentmwpb "github.com/NpoolPlatform/message/npool/account/mw/v1/payment"

	coininfopb "github.com/NpoolPlatform/message/npool/chain/mw/v1/coin"

	coininfocli "github.com/NpoolPlatform/chain-middleware/pkg/client/coin"
)

func GetAccount(ctx context.Context, id string) (*npool.Account, error) {
	info, err := paymentmwcli.GetAccount(ctx, id)
	if err != nil {
		return nil, err
	}

	coin, err := coininfocli.GetCoin(ctx, info.CoinTypeID)
	if err != nil {
		return nil, err
	}

	account := &npool.Account{
		ID:            info.ID,
		CoinTypeID:    info.CoinTypeID,
		CoinName:      coin.Name,
		CoinUnit:      coin.Unit,
		CoinEnv:       coin.ENV,
		CoinLogo:      coin.Logo,
		AccountID:     info.AccountID,
		Address:       info.Address,
		CollectingTID: info.CollectingTID,
		Active:        info.Active,
		Locked:        info.Locked,
		LockedBy:      info.LockedBy,
		Blocked:       info.Blocked,
		CreatedAt:     info.CreatedAt,
		AvailableAt:   info.AvailableAt,
		UpdatedAt:     info.UpdatedAt,
	}

	return account, nil
}

func GetAccounts(ctx context.Context, offset, limit int32) ([]*npool.Account, uint32, error) {
	infos, total, err := paymentmwcli.GetAccounts(ctx, &paymentmwpb.Conds{}, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	if len(infos) == 0 {
		return nil, total, nil
	}

	ofs := 0
	lim := 1000
	coins := []*coininfopb.Coin{}
	for {
		coinInfos, _, err := coininfocli.GetCoins(ctx, nil, int32(ofs), int32(lim))
		if err != nil {
			return nil, 0, err
		}
		if len(coinInfos) == 0 {
			break
		}
		coins = append(coins, coinInfos...)
		ofs += lim
	}

	coinMap := map[string]*coininfopb.Coin{}
	for _, coin := range coins {
		coinMap[coin.ID] = coin
	}

	accs := []*npool.Account{}

	for _, info := range infos {
		coin, ok := coinMap[info.CoinTypeID]
		if !ok {
			continue
		}

		accs = append(accs, &npool.Account{
			ID:            info.ID,
			CoinTypeID:    info.CoinTypeID,
			CoinName:      coin.Name,
			CoinUnit:      coin.Unit,
			CoinEnv:       coin.ENV,
			CoinLogo:      coin.Logo,
			AccountID:     info.AccountID,
			Address:       info.Address,
			CollectingTID: info.CollectingTID,
			Active:        info.Active,
			Locked:        info.Locked,
			LockedBy:      info.LockedBy,
			Blocked:       info.Blocked,
			CreatedAt:     info.CreatedAt,
			AvailableAt:   info.AvailableAt,
			UpdatedAt:     info.UpdatedAt,
		})
	}

	return accs, total, nil
}
