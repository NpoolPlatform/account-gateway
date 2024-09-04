package orderbenefit

import (
	"context"
	"fmt"

	orderbenefitmwcli "github.com/NpoolPlatform/account-middleware/pkg/client/orderbenefit"
	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	appcoinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/app/coin"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/orderbenefit"
	orderbenefitmwpb "github.com/NpoolPlatform/message/npool/account/mw/v1/orderbenefit"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	appcoinmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/app/coin"
)

type queryHandler struct {
	*Handler
	infos []*orderbenefitmwpb.Account
	coins map[string]*appcoinmwpb.Coin
	users map[string]*usermwpb.User
	accs  []*npool.Account
}

func (h *queryHandler) getUsers(ctx context.Context) error {
	ids := []string{}
	for _, info := range h.infos {
		ids = append(ids, info.UserID)
	}

	users, _, err := usermwcli.GetUsers(ctx, &usermwpb.Conds{
		EntIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: ids},
	}, 0, int32(len(ids)))
	if err != nil {
		return err
	}

	for _, u := range users {
		h.users[u.EntID] = u
	}

	return nil
}

func (h *queryHandler) getCoins(ctx context.Context) error {
	coinTypeIDs := []string{}
	for _, val := range h.infos {
		coinTypeIDs = append(coinTypeIDs, val.CoinTypeID)
	}

	coins, _, err := appcoinmwcli.GetCoins(ctx, &appcoinmwpb.Conds{
		AppID:       &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		CoinTypeIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: coinTypeIDs},
	}, 0, int32(len(coinTypeIDs)))
	if err != nil {
		return err
	}

	for _, coin := range coins {
		h.coins[coin.CoinTypeID] = coin
	}

	return nil
}

func (h *queryHandler) formalize() {
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
			EntID:            val.EntID,
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
			UsedFor:          val.UsedFor,
			OrderID:          val.OrderID,
			CreatedAt:        val.CreatedAt,
			PhoneNO:          userInfo.PhoneNO,
			EmailAddress:     userInfo.EmailAddress,
			Active:           val.Active,
			Blocked:          val.Blocked,
		})
	}
}

func (h *Handler) GetAccount(ctx context.Context) (*npool.Account, error) {
	info, err := orderbenefitmwcli.GetAccount(ctx, *h.EntID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf("invalid account")
	}
	return h.getAccountExt(ctx, info)
}

func (h *Handler) getAccountExt(ctx context.Context, info *orderbenefitmwpb.Account) (*npool.Account, error) {
	if info == nil {
		return nil, fmt.Errorf("invalid account")
	}
	handler := &queryHandler{
		Handler: h,
		infos:   []*orderbenefitmwpb.Account{info},
		coins:   map[string]*appcoinmwpb.Coin{},
		users:   map[string]*usermwpb.User{},
	}
	handler.AppID = &info.AppID
	if err := handler.getUsers(ctx); err != nil {
		return nil, err
	}
	if err := handler.getCoins(ctx); err != nil {
		return nil, err
	}
	if len(handler.users) == 0 {
		return nil, fmt.Errorf("invalid user")
	}
	if len(handler.coins) == 0 {
		return nil, fmt.Errorf("invalid coin")
	}
	handler.formalize()

	return handler.accs[0], nil
}

func (h *Handler) GetAccounts(ctx context.Context) ([]*npool.Account, uint32, error) {
	conds := &orderbenefitmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	}
	if h.UserID != nil {
		conds.UserID = &basetypes.StringVal{Op: cruder.EQ, Value: *h.UserID}
	}

	handler := &queryHandler{
		Handler: h,
		infos:   []*orderbenefitmwpb.Account{},
		users:   map[string]*usermwpb.User{},
		coins:   map[string]*appcoinmwpb.Coin{},
	}

	return handler.getAccounts(
		ctx,
		conds,
	)
}

func (h *queryHandler) getAccounts(ctx context.Context, conds *orderbenefitmwpb.Conds) ([]*npool.Account, uint32, error) {
	infos, total, err := orderbenefitmwcli.GetAccounts(ctx, conds, h.Offset, h.Limit)
	if err != nil {
		return nil, 0, err
	}
	if len(infos) == 0 {
		return nil, total, nil
	}

	h.infos = append(h.infos, infos...)

	if err := h.getUsers(ctx); err != nil {
		return nil, 0, err
	}
	if err := h.getCoins(ctx); err != nil {
		return nil, 0, err
	}
	h.formalize()

	return h.accs, total, nil
}
