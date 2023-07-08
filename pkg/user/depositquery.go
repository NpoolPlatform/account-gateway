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
	infos       []*depositmwpb.Account
	userIDs     []string
	coinTypeIDs []string
	coins       map[string]*appcoinmwpb.Coin
	users       map[string]*usermwpb.User
	accs        []*npool.Account
}

func (h *queryDepositHandler) getUsers(ctx context.Context) error {
	for _, info := range h.infos {
		h.userIDs = append(h.userIDs, info.UserID)
	}

	users, _, err := usermwcli.GetUsers(ctx, &usermwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		IDs:   &basetypes.StringSliceVal{Op: cruder.IN, Value: h.userIDs},
	}, 0, int32(len(h.userIDs)))
	if err != nil {
		return err
	}
	if len(users) == 0 {
		return fmt.Errorf("invalid user")
	}

	for _, u := range users {
		h.users[u.ID] = u
	}

	return nil
}

func (h *queryDepositHandler) getCoins(ctx context.Context) error {
	for _, val := range h.infos {
		h.coinTypeIDs = append(h.coinTypeIDs, val.CoinTypeID)
	}

	coins, _, err := appcoinmwcli.GetCoins(ctx, &appcoinmwpb.Conds{
		AppID: &basetypes.StringVal{
			Op:    cruder.EQ,
			Value: *h.AppID,
		},
		CoinTypeIDs: &basetypes.StringSliceVal{
			Op:    cruder.IN,
			Value: h.coinTypeIDs,
		},
	}, 0, int32(len(h.coinTypeIDs)))
	if err != nil {
		return err
	}
	if len(coins) == 0 {
		return fmt.Errorf("invalid coin")
	}

	for _, coin := range coins {
		h.coins[coin.CoinTypeID] = coin
	}

	return nil
}

func (h *queryDepositHandler) createAccount(ctx context.Context) error {
	sacc, err := sphinxproxycli.CreateAddress(ctx, h.coins[*h.CoinTypeID].CoinName)
	if err != nil {
		return err
	}
	if sacc == nil || sacc.Address == "" {
		return fmt.Errorf("fail create wallet")
	}

	info, err := depositcli.CreateAccount(ctx, &depositmwpb.AccountReq{
		AppID:      h.AppID,
		UserID:     h.UserID,
		CoinTypeID: h.CoinTypeID,
		Address:    &sacc.Address,
	})
	if err != nil {
		return err
	}
	h.infos = append(h.infos, info)

	return nil
}

func (h *queryDepositHandler) getBalance(ctx context.Context) error {
	bal, err := sphinxproxycli.GetBalance(ctx, &sphinxproxypb.GetBalanceRequest{
		Name:    h.coins[*h.CoinTypeID].CoinName,
		Address: h.infos[0].Address,
	})
	if err != nil {
		return err
	}
	if bal == nil {
		return fmt.Errorf("invalid address")
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

	handler := &queryDepositHandler{
		Handler:     h,
		infos:       []*depositmwpb.Account{},
		userIDs:     []string{*h.UserID},
		coinTypeIDs: []string{*h.CoinTypeID},
		coins:       map[string]*appcoinmwpb.Coin{},
		users:       map[string]*usermwpb.User{},
	}

	if err := handler.getUsers(ctx); err != nil {
		return nil, err
	}

	if err := handler.getCoins(ctx); err != nil {
		return nil, err
	}

	infos, _, err := depositcli.GetAccounts(ctx, &depositmwpb.Conds{
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
	if len(infos) > 0 {
		handler.infos = append(handler.infos, infos...)
	} else {
		handler.createAccount(ctx)
	}

	if err := handler.getBalance(ctx); err != nil {
		return nil, err
	}
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
		Handler:     h,
		infos:       infos,
		userIDs:     []string{},
		coinTypeIDs: []string{},
		coins:       map[string]*appcoinmwpb.Coin{},
		users:       map[string]*usermwpb.User{},
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
