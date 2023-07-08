//nolint:dupl
package user

import (
	"context"
	"fmt"

	depositcli "github.com/NpoolPlatform/account-middleware/pkg/client/deposit"
	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	appcoinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/app/coin"
	sphinxproxycli "github.com/NpoolPlatform/sphinx-proxy/pkg/client"

	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/user"
	depositmwpb "github.com/NpoolPlatform/message/npool/account/mw/v1/deposit"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	appcoinmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/app/coin"
	sphinxproxypb "github.com/NpoolPlatform/message/npool/sphinxproxy"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
)

type queryDepositHandler struct {
	*Handler
	infos []*depositmwpb.Account
	coins map[string]*appcoinmwpb.Coin
	users map[string]*usermwpb.User
	accs  []*npool.Account
}

func (h *queryDepositHandler) getUsers(ctx context.Context) error {
	ids := []string{}
	for _, info := range h.infos {
		ids = append(ids, info.UserID)
	}

	users, _, err := usermwcli.GetUsers(ctx, &usermwpb.Conds{
		IDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: ids},
	}, 0, int32(len(ids)))
	if err != nil {
		return err
	}

	for _, u := range users {
		h.users[u.ID] = u
	}

	return nil
}

func (h *queryDepositHandler) getCoins(ctx context.Context) error {
	coinTypeIDs := []string{}
	for _, val := range h.infos {
		coinTypeIDs = append(coinTypeIDs, val.CoinTypeID)
	}

	coins, _, err := appcoinmwcli.GetCoins(ctx, &appcoinmwpb.Conds{
		AppID: &basetypes.StringVal{
			Op:    cruder.EQ,
			Value: *h.AppID,
		},
		CoinTypeIDs: &basetypes.StringSliceVal{
			Op:    cruder.IN,
			Value: coinTypeIDs,
		},
	}, 0, int32(len(coinTypeIDs)))
	if err != nil {
		return err
	}

	for _, coin := range coins {
		h.coins[coin.CoinTypeID] = coin
	}

	return nil
}

func (h *queryDepositHandler) formalize() {
	for _, val := range h.infos {
		userInfo, ok := h.users[val.UserID]
		if !ok {
			continue
		}
		coin, ok := h.coins[val.CoinTypeID]
		if !ok {
			continue
		}

		h.accs = append(h.accs, &npool.Account{
			ID:               val.ID,
			AppID:            val.AppID,
			UserID:           val.UserID,
			CoinTypeID:       val.CoinTypeID,
			CoinName:         coin.Name,
			CoinDisplayNames: coin.DisplayNames,
			CoinUnit:         coin.Unit,
			CoinEnv:          coin.ENV,
			CoinLogo:         coin.Logo,
			AccountID:        val.AccountID,
			Address:          val.Address,
			CreatedAt:        val.CreatedAt,
			PhoneNO:          userInfo.PhoneNO,
			EmailAddress:     userInfo.EmailAddress,
		})
	}
}

func (h *Handler) GetDepositAccount(ctx context.Context) (*npool.Account, error) { //nolint
	if h.AppID == nil {
		return nil, fmt.Errorf("invaild appid")
	}
	if h.UserID == nil {
		return nil, fmt.Errorf("invaild userID")
	}
	if h.CoinTypeID == nil {
		return nil, fmt.Errorf("invaild coinTypeID")
	}

	user, err := usermwcli.GetUser(ctx, *h.AppID, *h.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("invalid user")
	}
	if user.AppID != *h.AppID {
		return nil, fmt.Errorf("permission denied")
	}

	coin, err := appcoinmwcli.GetCoinOnly(ctx, &appcoinmwpb.Conds{
		AppID:      &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		CoinTypeID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.CoinTypeID},
	})
	if err != nil {
		return nil, err
	}
	if coin == nil {
		return nil, fmt.Errorf("invalid coin")
	}

	handler := &queryDepositHandler{
		Handler: h,
		infos:   []*depositmwpb.Account{},
		coins:   map[string]*appcoinmwpb.Coin{},
		users:   map[string]*usermwpb.User{},
	}
	handler.users[user.ID] = user
	handler.coins[coin.CoinTypeID] = coin

	accs, _, err := depositcli.GetAccounts(ctx, &depositmwpb.Conds{
		AppID:      &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		UserID:     &basetypes.StringVal{Op: cruder.EQ, Value: *h.UserID},
		CoinTypeID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.CoinTypeID},
		Active:     &basetypes.BoolVal{Op: cruder.EQ, Value: true},
		Locked:     &basetypes.BoolVal{Op: cruder.EQ, Value: false},
		Blocked:    &basetypes.BoolVal{Op: cruder.EQ, Value: false},
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
		handler.infos = append(handler.infos, acc)
		handler.formalize()

		return handler.accs[0], nil
	}

	sacc, err := sphinxproxycli.CreateAddress(ctx, coin.Name)
	if err != nil {
		return nil, err
	}
	if sacc == nil || sacc.Address == "" {
		return nil, fmt.Errorf("fail create wallet")
	}

	acc, err := depositcli.CreateAccount(ctx, &depositmwpb.AccountReq{
		AppID:      h.AppID,
		UserID:     h.UserID,
		CoinTypeID: h.CoinTypeID,
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
	handler.infos = append(handler.infos, acc)
	handler.formalize()

	return handler.accs[0], nil
}

func (h *Handler) GetDepositAccounts(ctx context.Context) ([]*npool.Account, uint32, error) {
	if h.AppID == nil {
		return nil, 0, fmt.Errorf("invaild appid")
	}

	infos, total, err := depositcli.GetAccounts(ctx, &depositmwpb.Conds{
		AppID: &basetypes.StringVal{
			Op:    cruder.EQ,
			Value: *h.AppID,
		},
	}, h.Offset, h.Limit)
	if err != nil {
		return nil, 0, err
	}

	if len(infos) == 0 {
		return nil, 0, nil
	}

	handler := &queryDepositHandler{
		Handler: h,
		infos:   infos,
		coins:   map[string]*appcoinmwpb.Coin{},
		users:   map[string]*usermwpb.User{},
	}
	if err := handler.getUsers(ctx); err != nil {
		return nil, 0, err
	}
	if err := handler.getCoins(ctx); err != nil {
		return nil, 0, err
	}
	handler.formalize()

	return handler.accs, total, nil
}
