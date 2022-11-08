package user

import (
	"context"
	"fmt"

	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/user"

	useraccmwcli "github.com/NpoolPlatform/account-middleware/pkg/client/user"
	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	coininfocli "github.com/NpoolPlatform/sphinx-coininfo/pkg/client"

	commonpb "github.com/NpoolPlatform/message/npool"
	useraccmwpb "github.com/NpoolPlatform/message/npool/account/mw/v1/user"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	coininfopb "github.com/NpoolPlatform/message/npool/coininfo"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
)

func GetAccount(ctx context.Context, id string) (*npool.Account, error) {
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

	fmt.Println("info: ", info)
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
		Active:       info.Active,
		Blocked:      info.Blocked,
		Labels:       info.Labels,
	}
	return acc, nil
}

func GetAccounts(ctx context.Context, appID, userID string, offset, limit int32) ([]*npool.Account, uint32, error) {
	return getAccounts(
		ctx,
		&useraccmwpb.Conds{
			AppID: &commonpb.StringVal{
				Op:    cruder.EQ,
				Value: appID,
			},
			UserID: &commonpb.StringVal{
				Op:    cruder.EQ,
				Value: userID,
			},
		},
		offset,
		limit,
	)
}

func GetAppAccounts(ctx context.Context, appID string, offset, limit int32) ([]*npool.Account, uint32, error) {
	return getAccounts(
		ctx,
		&useraccmwpb.Conds{
			AppID: &commonpb.StringVal{
				Op:    cruder.EQ,
				Value: appID,
			},
		},
		offset,
		limit,
	)
}

func getAccounts(ctx context.Context, conds *useraccmwpb.Conds, offset, limit int32) ([]*npool.Account, uint32, error) {
	infos, total, err := useraccmwcli.GetAccounts(ctx, conds, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	if len(infos) == 0 {
		return nil, total, nil
	}


	for _, info := range infos {
		fmt.Println("CreatedAt: ", info.CreatedAt)
	}

	ids := []string{}
	for _, info := range infos {
		ids = append(ids, info.UserID)
	}

	users, _, err := usermwcli.GetManyUsers(ctx, ids)
	if err != nil {
		return nil, 0, err
	}

	userMap := map[string]*usermwpb.User{}
	for _, u := range users {
		userMap[u.ID] = u
	}

	coins, err := coininfocli.GetCoinInfos(ctx, cruder.NewFilterConds())
	if err != nil {
		return nil, 0, err
	}

	coinMap := map[string]*coininfopb.CoinInfo{}
	for _, coin := range coins {
		coinMap[coin.ID] = coin
	}

	accs := []*npool.Account{}
	for _, info := range infos {
		u, ok := userMap[info.UserID]
		if !ok {
			continue
		}

		coin, ok := coinMap[info.CoinTypeID]
		if !ok {
			continue
		}

		accs = append(accs, &npool.Account{
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
			Active:       info.Active,
			Blocked:      info.Blocked,
			Labels:       info.Labels,
		})
	}

	return accs, total, nil
}
