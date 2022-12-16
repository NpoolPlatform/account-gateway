package user

import (
	"context"
	"fmt"

	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/user"

	depositcli "github.com/NpoolPlatform/account-middleware/pkg/client/deposit"
	depositpb "github.com/NpoolPlatform/message/npool/account/mw/v1/deposit"

	usercli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"

	appcoininfocli "github.com/NpoolPlatform/chain-middleware/pkg/client/appcoin"
	coininfocli "github.com/NpoolPlatform/chain-middleware/pkg/client/coin"
	sphinxproxycli "github.com/NpoolPlatform/sphinx-proxy/pkg/client"

	appcoinpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/appcoin"
	sphinxproxypb "github.com/NpoolPlatform/message/npool/sphinxproxy"

	commonpb "github.com/NpoolPlatform/message/npool"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"

	appusermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
)

func GetDepositAccount(ctx context.Context, appID, userID, coinTypeID string) (*npool.Account, error) { //nolint
	user, err := usercli.GetUser(ctx, appID, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("invalid user")
	}
	if user.AppID != appID {
		return nil, fmt.Errorf("permission denied")
	}

	coin, err := coininfocli.GetCoin(ctx, coinTypeID)
	if err != nil {
		return nil, err
	}
	if coin == nil {
		return nil, fmt.Errorf("invalid coin")
	}

	accs, _, err := depositcli.GetAccounts(ctx, &depositpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: appID,
		},
		UserID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: userID,
		},
		CoinTypeID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: coinTypeID,
		},
		Active: &commonpb.BoolVal{
			Op:    cruder.EQ,
			Value: true,
		},
		Locked: &commonpb.BoolVal{
			Op:    cruder.EQ,
			Value: false,
		},
		Blocked: &commonpb.BoolVal{
			Op:    cruder.EQ,
			Value: false,
		},
	}, 0, 1)
	if err != nil {
		return nil, err
	}

	if len(accs) > 0 {
		acc := accs[0]

		bal, err := sphinxproxycli.GetBalance(ctx, &sphinxproxypb.GetBalanceRequest{
			Name:    coin.Name,
			Address: acc.Address,
		})
		if err != nil {
			return nil, err
		}
		if bal == nil {
			return nil, fmt.Errorf("invalid address")
		}

		return &npool.Account{
			ID:         acc.ID,
			AppID:      acc.AppID,
			UserID:     acc.UserID,
			CoinTypeID: acc.CoinTypeID,
			CoinName:   coin.Name,
			CoinUnit:   coin.Unit,
			CoinEnv:    coin.ENV,
			CoinLogo:   coin.Logo,
			AccountID:  acc.AccountID,
			Address:    acc.Address,
			CreatedAt:  acc.CreatedAt,
		}, nil
	}

	sacc, err := sphinxproxycli.CreateAddress(ctx, coin.Name)
	if err != nil {
		return nil, err
	}
	if sacc == nil || sacc.Address == "" {
		return nil, fmt.Errorf("fail create wallet")
	}

	acc, err := depositcli.CreateAccount(ctx, &depositpb.AccountReq{
		AppID:      &appID,
		UserID:     &userID,
		CoinTypeID: &coinTypeID,
		Address:    &sacc.Address,
	})
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

	return &npool.Account{
		ID:         acc.ID,
		AppID:      acc.AppID,
		UserID:     acc.UserID,
		CoinTypeID: acc.CoinTypeID,
		CoinName:   coin.Name,
		CoinUnit:   coin.Unit,
		CoinEnv:    coin.ENV,
		CoinLogo:   coin.Logo,
		AccountID:  acc.AccountID,
		Address:    acc.Address,
		CreatedAt:  acc.CreatedAt,
	}, nil
}

//nolint
func GetDepositAccounts(ctx context.Context, appID string, offset, limit int32) ([]*npool.Account, uint32, error) {
	accs, total, err := depositcli.GetAccounts(ctx, &depositpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: appID,
		},
	}, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	if len(accs) == 0 {
		return nil, 0, nil
	}

	userIDs := []string{}
	for _, info := range accs {
		userIDs = append(userIDs, info.UserID)
	}

	users, _, err := usercli.GetManyUsers(ctx, userIDs)
	if err != nil {
		return nil, 0, fmt.Errorf("fail get users: %v", err)
	}

	userMap := map[string]*appusermwpb.User{}
	for _, user := range users {
		userMap[user.ID] = user
	}

	coinTypeIDs := []string{}
	for _, val := range accs {
		coinTypeIDs = append(coinTypeIDs, val.CoinTypeID)
	}

	coins, _, err := appcoininfocli.GetCoins(ctx, &appcoinpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: appID,
		},
		CoinTypeIDs: &commonpb.StringSliceVal{
			Op:    cruder.IN,
			Value: coinTypeIDs,
		},
	}, 0, int32(len(coinTypeIDs)))
	if err != nil {
		return nil, 0, err
	}

	coinMap := map[string]*appcoinpb.Coin{}
	for _, coin := range coins {
		coinMap[coin.CoinTypeID] = coin
	}

	infos := []*npool.Account{}
	for _, acc := range accs {
		coin, ok := coinMap[acc.CoinTypeID]
		if !ok {
			continue
		}
		user, ok := userMap[acc.UserID]
		if !ok {
			continue
		}
		infos = append(infos, &npool.Account{
			ID:           acc.ID,
			AppID:        acc.AppID,
			UserID:       acc.UserID,
			CoinTypeID:   acc.CoinTypeID,
			CoinName:     coin.Name,
			CoinUnit:     coin.Unit,
			CoinEnv:      coin.ENV,
			CoinLogo:     coin.Logo,
			AccountID:    acc.AccountID,
			Address:      acc.Address,
			CreatedAt:    acc.CreatedAt,
			PhoneNO:      user.PhoneNO,
			EmailAddress: user.EmailAddress,
		})
	}
	return infos, total, nil
}

//nolint
func GetAppDepositAccounts(ctx context.Context, appID string, offset, limit int32) ([]*npool.Account, uint32, error) {
	accs, total, err := depositcli.GetAccounts(ctx, &depositpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: appID,
		},
	}, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	if len(accs) == 0 {
		return nil, 0, nil
	}

	userIDs := []string{}
	for _, info := range accs {
		userIDs = append(userIDs, info.UserID)
	}

	users, _, err := usercli.GetManyUsers(ctx, userIDs)
	if err != nil {
		return nil, 0, fmt.Errorf("fail get users: %v", err)
	}

	userMap := map[string]*appusermwpb.User{}
	for _, user := range users {
		userMap[user.ID] = user
	}

	coinTypeIDs := []string{}
	for _, val := range accs {
		coinTypeIDs = append(coinTypeIDs, val.CoinTypeID)
	}

	coins, _, err := appcoininfocli.GetCoins(ctx, &appcoinpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: appID,
		},
		CoinTypeIDs: &commonpb.StringSliceVal{
			Op:    cruder.IN,
			Value: coinTypeIDs,
		},
	}, 0, int32(len(coinTypeIDs)))
	if err != nil {
		return nil, 0, err
	}

	coinMap := map[string]*appcoinpb.Coin{}
	for _, coin := range coins {
		coinMap[coin.CoinTypeID] = coin
	}

	infos := []*npool.Account{}
	for _, acc := range accs {
		coin, ok := coinMap[acc.CoinTypeID]
		if !ok {
			continue
		}
		user, ok := userMap[acc.UserID]
		if !ok {
			continue
		}
		infos = append(infos, &npool.Account{
			ID:           acc.ID,
			AppID:        acc.AppID,
			UserID:       acc.UserID,
			CoinTypeID:   acc.CoinTypeID,
			CoinName:     coin.Name,
			CoinUnit:     coin.Unit,
			CoinEnv:      coin.ENV,
			CoinLogo:     coin.Logo,
			AccountID:    acc.AccountID,
			Address:      acc.Address,
			CreatedAt:    acc.CreatedAt,
			PhoneNO:      user.PhoneNO,
			EmailAddress: user.EmailAddress,
		})
	}
	return infos, total, nil
}
