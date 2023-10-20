package payment

import (
	"context"
	"fmt"

	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/payment"

	paymentmwcli "github.com/NpoolPlatform/account-middleware/pkg/client/payment"
	paymentmwpb "github.com/NpoolPlatform/message/npool/account/mw/v1/payment"

	coinmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/coin"

	coinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/coin"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

type queryHandler struct {
	*Handler
	infos []*paymentmwpb.Account
	accs  []*npool.Account
	coins map[string]*coinmwpb.Coin
}

func (h *queryHandler) getCoins(ctx context.Context) error {
	coinTypeIDs := []string{}

	for _, info := range h.infos {
		coinTypeIDs = append(coinTypeIDs, info.CoinTypeID)
	}
	coins, _, err := coinmwcli.GetCoins(
		ctx,
		&coinmwpb.Conds{
			EntIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: coinTypeIDs},
		},
		0,
		int32(len(coinTypeIDs)),
	)
	if err != nil {
		return err
	}

	for _, coin := range coins {
		h.coins[coin.EntID] = coin
	}

	return nil
}

func (h *queryHandler) formalize() {
	for _, info := range h.infos {
		coin, ok := h.coins[info.CoinTypeID]
		if !ok {
			continue
		}

		h.accs = append(h.accs, &npool.Account{
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
}

func (h *Handler) GetAccount(ctx context.Context) (*npool.Account, error) {
	if h.ID == nil {
		return nil, fmt.Errorf("invalid id")
	}

	info, err := paymentmwcli.GetAccount(ctx, *h.ID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}

	handler := &queryHandler{
		Handler: h,
		infos:   []*paymentmwpb.Account{info},
		coins:   map[string]*coinmwpb.Coin{},
	}
	if err := handler.getCoins(ctx); err != nil {
		return nil, err
	}
	handler.formalize()

	if len(handler.accs) == 0 {
		return nil, nil
	}

	return handler.accs[0], nil
}

func (h *Handler) GetAccounts(ctx context.Context) ([]*npool.Account, uint32, error) {
	infos, total, err := paymentmwcli.GetAccounts(
		ctx,
		&paymentmwpb.Conds{},
		h.Offset,
		h.Limit,
	)
	if err != nil {
		return nil, 0, err
	}
	if len(infos) == 0 {
		return nil, total, nil
	}

	handler := &queryHandler{
		Handler: h,
		infos:   infos,
		coins:   map[string]*coinmwpb.Coin{},
	}
	if err := handler.getCoins(ctx); err != nil {
		return nil, 0, err
	}
	handler.formalize()

	return handler.accs, total, nil
}
