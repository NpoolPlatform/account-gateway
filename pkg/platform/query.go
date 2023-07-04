package platform

import (
	"context"
	"fmt"

	pltfmwcli "github.com/NpoolPlatform/account-middleware/pkg/client/platform"
	coinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/coin"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/platform"
	pltfmwpb "github.com/NpoolPlatform/message/npool/account/mw/v1/platform"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	coinmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/coin"
)

type queryHandler struct {
	*Handler
	infos []*pltfmwpb.Account
	accs  []*npool.Account
	// GoodID -> Coin
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
			IDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: coinTypeIDs},
		},
		0,
		int32(len(coinTypeIDs)),
	)
	if err != nil {
		return err
	}

	for _, coin := range coins {
		h.coins[coin.ID] = coin
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
			ID:         info.ID,
			CoinTypeID: info.CoinTypeID,
			CoinName:   coin.Name,
			CoinUnit:   coin.Unit,
			CoinEnv:    coin.ENV,
			CoinLogo:   coin.Logo,
			UsedFor:    info.UsedFor,
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
}

func (h *Handler) GetAccount(ctx context.Context) (*npool.Account, error) {
	if h.ID == nil {
		return nil, fmt.Errorf("invalid id")
	}

	info, err := pltfmwcli.GetAccount(ctx, *h.ID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}

	handler := &queryHandler{
		Handler: h,
		infos:   []*pltfmwpb.Account{info},
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
	infos, total, err := pltfmwcli.GetAccounts(
		ctx,
		&pltfmwpb.Conds{},
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
