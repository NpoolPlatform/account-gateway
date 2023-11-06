package goodbenefit

import (
	"context"

	gbmwcli "github.com/NpoolPlatform/account-middleware/pkg/client/goodbenefit"
	coinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/coin"
	goodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/good"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/goodbenefit"
	gbmwpb "github.com/NpoolPlatform/message/npool/account/mw/v1/goodbenefit"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	coinmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/coin"
	goodmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/good"
)

type queryHandler struct {
	*Handler
	infos []*gbmwpb.Account
	accs  []*npool.Account
	// GoodID -> Coin
	coins map[string]*coinmwpb.Coin
	goods map[string]*goodmwpb.Good
}

func (h *queryHandler) getGoods(ctx context.Context) error {
	goodIDs := []string{}
	for _, info := range h.infos {
		goodIDs = append(goodIDs, info.GoodID)
	}
	goods, _, err := goodmwcli.GetGoods(ctx, &goodmwpb.Conds{
		IDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: goodIDs},
	}, 0, int32(len(goodIDs)))
	if err != nil {
		return err
	}
	for _, good := range goods {
		h.goods[good.ID] = good
	}
	return nil
}

func (h *queryHandler) getCoins(ctx context.Context) error {
	coinTypeIDs := []string{}
	coinGoodIDs := map[string][]string{}

	for _, good := range h.goods {
		coinTypeIDs = append(coinTypeIDs, good.CoinTypeID)
		coinGoodIDs[good.CoinTypeID] = append(coinGoodIDs[good.CoinTypeID], good.ID)
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
		goodIDs, ok := coinGoodIDs[coin.EntID]
		if !ok {
			continue
		}
		for _, goodID := range goodIDs {
			h.coins[goodID] = coin
		}
	}

	return nil
}

func (h *queryHandler) formalize() {
	for _, info := range h.infos {
		good, ok := h.goods[info.GoodID]
		if !ok {
			continue
		}

		coin, ok := h.coins[info.GoodID]
		if !ok {
			continue
		}

		h.accs = append(h.accs, &npool.Account{
			ID:         info.ID,
			EntID:      info.EntID,
			GoodID:     info.GoodID,
			GoodName:   good.Title,
			GoodUnit:   good.Unit,
			CoinTypeID: info.CoinTypeID,
			CoinName:   coin.Name,
			CoinUnit:   coin.Unit,
			CoinEnv:    coin.ENV,
			CoinLogo:   coin.Logo,
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
	info, err := gbmwcli.GetAccount(ctx, *h.EntID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}

	handler := &queryHandler{
		Handler: h,
		infos:   []*gbmwpb.Account{info},
		coins:   map[string]*coinmwpb.Coin{},
		goods:   map[string]*goodmwpb.Good{},
	}
	if err := handler.getGoods(ctx); err != nil {
		return nil, err
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
	infos, total, err := gbmwcli.GetAccounts(
		ctx,
		&gbmwpb.Conds{},
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
		goods:   map[string]*goodmwpb.Good{},
	}
	if err := handler.getGoods(ctx); err != nil {
		return nil, 0, err
	}
	if err := handler.getCoins(ctx); err != nil {
		return nil, 0, err
	}
	handler.formalize()

	return handler.accs, total, nil
}
