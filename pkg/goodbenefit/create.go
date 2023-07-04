package goodbenefit

import (
	"context"
	"fmt"

	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/goodbenefit"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	gbmwcli "github.com/NpoolPlatform/account-middleware/pkg/client/goodbenefit"
	gbmwpb "github.com/NpoolPlatform/message/npool/account/mw/v1/goodbenefit"

	coinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/coin"
	goodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/good"
	sphinxproxypb "github.com/NpoolPlatform/message/npool/sphinxproxy"
	sphinxproxycli "github.com/NpoolPlatform/sphinx-proxy/pkg/client"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
)

type createHandler struct {
	*Handler
	goodCoinTypeID      *string
	goodCoinName        *string
	checkAddressBalance bool
	backup              bool
	address             *string
}

func (h *createHandler) getCoinTypeID(ctx context.Context) error {
	if h.GoodID == nil {
		return fmt.Errorf("invalid goodid")
	}

	good, err := goodmwcli.GetGood(ctx, *h.GoodID)
	if err != nil {
		return err
	}
	if good == nil {
		return fmt.Errorf("invalid good")
	}

	h.goodCoinTypeID = &good.CoinTypeID
	return nil
}

func (h *createHandler) getCoinName(ctx context.Context) error {
	if h.CoinTypeID == nil {
		return fmt.Errorf("invalid goodcointypeid")
	}

	coin, err := coinmwcli.GetCoin(ctx, *h.CoinTypeID)
	if err != nil {
		return err
	}
	if coin == nil {
		return fmt.Errorf("invalid coin")
	}

	h.goodCoinName = &coin.Name
	h.checkAddressBalance = coin.CheckNewAddressBalance

	return nil
}

func (h *createHandler) checkBackup(ctx context.Context) error {
	exist, err := gbmwcli.ExistAccountConds(ctx, &gbmwpb.Conds{
		GoodID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.GoodID},
		Backup: &basetypes.BoolVal{Op: cruder.EQ, Value: false},
	})
	if err != nil {
		return err
	}
	h.backup = exist
	return nil
}

func (h *createHandler) createAddress(ctx context.Context) error {
	if h.goodCoinName == nil {
		return fmt.Errorf("invalid goodcoinname")
	}

	acc, err := sphinxproxycli.CreateAddress(ctx, *h.goodCoinName)
	if err != nil {
		return err
	}
	if acc == nil {
		return fmt.Errorf("fail create address")
	}

	h.address = &acc.Address

	if !h.checkAddressBalance {
		return nil
	}

	bal, err := sphinxproxycli.GetBalance(ctx, &sphinxproxypb.GetBalanceRequest{
		Name:    *h.goodCoinName,
		Address: acc.Address,
	})
	if err != nil {
		return err
	}
	if bal == nil {
		return fmt.Errorf("invalid address")
	}

	return nil
}

func (h *createHandler) createAccount(ctx context.Context) error {
	acc, err := gbmwcli.CreateAccount(ctx, &gbmwpb.AccountReq{
		GoodID:     h.GoodID,
		CoinTypeID: h.goodCoinTypeID,
		Address:    h.address,
		Backup:     &h.backup,
	})
	if err != nil {
		return err
	}
	if acc == nil {
		return fmt.Errorf("fail create account")
	}

	h.ID = &acc.ID

	return nil
}

func (h *Handler) CreateAccount(ctx context.Context) (*npool.Account, error) {
	handler := &createHandler{
		Handler: h,
	}

	if err := handler.getCoinTypeID(ctx); err != nil {
		return nil, err
	}
	if err := handler.getCoinName(ctx); err != nil {
		return nil, err
	}
	if err := handler.checkBackup(ctx); err != nil {
		return nil, err
	}
	if err := handler.createAddress(ctx); err != nil {
		return nil, err
	}
	if err := handler.createAccount(ctx); err != nil {
		return nil, err
	}

	return h.GetAccount(ctx)
}
