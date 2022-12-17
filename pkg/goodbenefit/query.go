package goodbenefit

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/goodbenefit"

	gbmwcli "github.com/NpoolPlatform/account-middleware/pkg/client/goodbenefit"

	gbmwpb "github.com/NpoolPlatform/message/npool/account/mw/v1/goodbenefit"

	goodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/good"
	goodmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/good"

	coininfopb "github.com/NpoolPlatform/message/npool/chain/mw/v1/coin"

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
		Address:    info.Address,
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
	infos, total, err := gbmwcli.GetAccounts(ctx, &gbmwpb.Conds{}, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	if len(infos) == 0 {
		return nil, total, nil
	}

	ids := []string{}
	for _, info := range infos {
		ids = append(ids, info.GoodID)
	}

	goods, _, err := goodmwcli.GetManyGoods(ctx, ids, 0, int32(len(ids)))
	if err != nil {
		return nil, 0, err
	}

	goodMap := map[string]*goodmwpb.Good{}
	for _, good := range goods {
		goodMap[good.ID] = good
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
		good, ok := goodMap[info.GoodID]
		if !ok {
			continue
		}

		coin, ok := coinMap[info.CoinTypeID]
		if !ok {
			continue
		}

		accs = append(accs, &npool.Account{
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
			Address:    info.Address,
			Active:     info.Active,
			Locked:     info.Locked,
			LockedBy:   info.LockedBy,
			Blocked:    info.Blocked,
			CreatedAt:  info.CreatedAt,
			UpdatedAt:  info.UpdatedAt,
		})
	}

	return accs, total, nil
}
