package platform

import (
	"context"
	"fmt"

	commonpb "github.com/NpoolPlatform/message/npool"
	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/platform"

	gbmwcli "github.com/NpoolPlatform/account-middleware/pkg/client/platform"
	gbmwpb "github.com/NpoolPlatform/message/npool/account/mw/v1/platform"

	accountmgrpb "github.com/NpoolPlatform/message/npool/account/mgr/v1/account"

	sphinxproxypb "github.com/NpoolPlatform/message/npool/sphinxproxy"
	coininfocli "github.com/NpoolPlatform/sphinx-coininfo/pkg/client"
	sphinxproxycli "github.com/NpoolPlatform/sphinx-proxy/pkg/client"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
)

func CreateAccount(
	ctx context.Context,
	coinTypeID string,
	address *string,
	usedFor accountmgrpb.AccountUsedFor,
) (
	*npool.Account, error,
) {
	coin, err := coininfocli.GetCoinInfo(ctx, coinTypeID)
	if err != nil {
		return nil, err
	}
	if coin == nil {
		return nil, fmt.Errorf("invalid coin")
	}

	backup := false
	const accountNumber = 100

	accounts, _, err := gbmwcli.GetAccounts(ctx, &gbmwpb.Conds{
		CoinTypeID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: coinTypeID,
		},
		UsedFor: &commonpb.Int32Val{
			Op:    cruder.EQ,
			Value: int32(usedFor),
		},
	}, 0, accountNumber)
	if err != nil {
		return nil, err
	}
	if address != nil {
		for _, acc := range accounts {
			if acc.Address == *address {
				return GetAccount(ctx, acc.ID)
			}
		}
	}

	for _, acc := range accounts {
		if acc.Active && !acc.Blocked && !acc.Backup {
			backup = true
			break
		}
	}

	targetAddress := ""
	if address != nil {
		targetAddress = *address
	} else {
		sacc, err := sphinxproxycli.CreateAddress(ctx, coin.Name)
		if err != nil {
			return nil, err
		}
		if sacc == nil {
			return nil, fmt.Errorf("fail create address")
		}
		targetAddress = sacc.Address
	}

	bal, err := sphinxproxycli.GetBalance(ctx, &sphinxproxypb.GetBalanceRequest{
		Name:    coin.Name,
		Address: targetAddress,
	})
	if err != nil {
		return nil, err
	}
	if bal == nil {
		return nil, fmt.Errorf("invalid address")
	}

	acc, err := gbmwcli.CreateAccount(ctx, &gbmwpb.AccountReq{
		CoinTypeID: &coinTypeID,
		UsedFor:    &usedFor,
		Address:    &targetAddress,
		Backup:     &backup,
	})
	if err != nil {
		return nil, err
	}
	if acc == nil {
		return nil, fmt.Errorf("fail create account")
	}

	return GetAccount(ctx, acc.ID)
}
