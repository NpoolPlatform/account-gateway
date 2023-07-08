//nolint
package user

import (
	"context"
	"fmt"

	useraccmwcli "github.com/NpoolPlatform/account-middleware/pkg/client/user"
	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	appcoinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/app/coin"

	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/user"
	useraccmwpb "github.com/NpoolPlatform/message/npool/account/mw/v1/user"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	appcoinmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/app/coin"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
)

type queryHandler struct {
	*Handler
	infos []*useraccmwpb.Account
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
			CreatedAt:        val.CreatedAt,
			PhoneNO:          userInfo.PhoneNO,
			EmailAddress:     userInfo.EmailAddress,
			Active:           val.Active,
			Blocked:          val.Blocked,
			Labels:           val.Labels,
			Memo:             val.Memo,
		})
	}
}

func (h *Handler) GetAccount(ctx context.Context) (*npool.Account, error) {
	if h.ID == nil {
		return nil, fmt.Errorf("invalid id")
	}

	info, err := useraccmwcli.GetAccount(ctx, *h.ID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf("invalid account")
	}
	handler := &queryHandler{
		Handler: h,
		infos:   []*useraccmwpb.Account{info},
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
	if h.AppID == nil {
		return nil, 0, fmt.Errorf("invalid appID")
	}
	if h.UserID == nil {
		return nil, 0, fmt.Errorf("invalid userID")
	}
	if h.UsedFor == nil {
		return nil, 0, fmt.Errorf("invalid usedFor")
	}
	handler := &queryHandler{
		Handler: h,
		users:   map[string]*usermwpb.User{},
		coins:   map[string]*appcoinmwpb.Coin{},
	}

	return handler.getAccounts(
		ctx,
		&useraccmwpb.Conds{
			AppID:   &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
			UserID:  &basetypes.StringVal{Op: cruder.EQ, Value: *h.UserID},
			UsedFor: &basetypes.Uint32Val{Op: cruder.EQ, Value: uint32(*h.UsedFor)},
		},
	)
}

func (h *Handler) GetAppAccounts(ctx context.Context) ([]*npool.Account, uint32, error) {
	if h.AppID == nil {
		return nil, 0, fmt.Errorf("invalid appID")
	}
	handler := &queryHandler{
		Handler: h,
		users:   map[string]*usermwpb.User{},
		coins:   map[string]*appcoinmwpb.Coin{},
	}

	return handler.getAccounts(
		ctx,
		&useraccmwpb.Conds{
			AppID: &basetypes.StringVal{
				Op:    cruder.EQ,
				Value: *h.AppID,
			},
		},
	)
}

func (h *queryHandler) getAccounts(ctx context.Context, conds *useraccmwpb.Conds) ([]*npool.Account, uint32, error) {
	infos, total, err := useraccmwcli.GetAccounts(ctx, conds, h.Offset, h.Limit)
	if err != nil {
		return nil, 0, err
	}
	if len(infos) == 0 {
		return nil, total, nil
	}

	if err := h.getUsers(ctx); err != nil {
		return nil, 0, err
	}
	if err := h.getCoins(ctx); err != nil {
		return nil, 0, err
	}
	h.formalize()

	return h.accs, total, nil
}
