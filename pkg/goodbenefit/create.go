package goodbenefit

import (
	"context"
	"fmt"

	commonpb "github.com/NpoolPlatform/message/npool"
	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/goodbenefit"

	gbmwcli "github.com/NpoolPlatform/account-middleware/pkg/client/goodbenefit"
	gbmwpb "github.com/NpoolPlatform/message/npool/account/mw/v1/goodbenefit"

	goodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/good"

	coininfocli "github.com/NpoolPlatform/chain-middleware/pkg/client/coin"
	sphinxproxypb "github.com/NpoolPlatform/message/npool/sphinxproxy"
	sphinxproxycli "github.com/NpoolPlatform/sphinx-proxy/pkg/client"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
)

func CreateAccount(ctx context.Context, goodID string) (*npool.Account, error) {
	good, err := goodmwcli.GetGood(ctx, goodID)
	if err != nil {
		return nil, err
	}

	coin, err := coininfocli.GetCoin(ctx, good.CoinTypeID)
	if err != nil {
		return nil, err
	}

	backup := false
	const accountNumber = 100

	accounts, err := gbmwcli.GetAccounts(ctx, &gbmwpb.Conds{
		GoodID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: goodID,
		},
	}, 0, accountNumber)
	if err != nil {
		return nil, err
	}

	for _, acc := range accounts {
		if acc.Active && !acc.Blocked && !acc.Backup {
			backup = true
			break
		}
	}

	sacc, err := sphinxproxycli.CreateAddress(ctx, coin.Name)
	if err != nil {
		return nil, err
	}

	bal, err := sphinxproxycli.GetBalance(ctx, &sphinxproxypb.GetBalanceRequest{
		Name:    coin.Name,
		Address: sacc.Address,
	})
	if err != nil {
		return nil, err
	}
	if bal == nil {
		return nil, fmt.Errorf("invalid address")
	}

	acc, err := gbmwcli.CreateAccount(ctx, &gbmwpb.AccountReq{
		GoodID:     &goodID,
		CoinTypeID: &good.CoinTypeID,
		Address:    &sacc.Address,
		Backup:     &backup,
	})
	if err != nil {
		return nil, err
	}

	return GetAccount(ctx, acc.ID)
}
