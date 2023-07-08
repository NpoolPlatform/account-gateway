package platform

import (
	"context"
	"fmt"

	pltfmwcli "github.com/NpoolPlatform/account-middleware/pkg/client/platform"
	coinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/coin"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/platform"
	pltfmwpb "github.com/NpoolPlatform/message/npool/account/mw/v1/platform"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	sphinxproxypb "github.com/NpoolPlatform/message/npool/sphinxproxy"
	sphinxproxycli "github.com/NpoolPlatform/sphinx-proxy/pkg/client"
)

type createHandler struct {
	*Handler
	coinName            *string
	backup              bool
	checkAddressBalance bool
}

func (h *createHandler) checkAddress(ctx context.Context) error {
	if h.UsedFor == nil {
		return fmt.Errorf("invalid usedfor")
	}
	switch *h.UsedFor {
	case basetypes.AccountUsedFor_UserBenefitHot:
		fallthrough // nolint
	case basetypes.AccountUsedFor_GasProvider:
		if h.Address != nil {
			return fmt.Errorf("invalid address")
		}
	case basetypes.AccountUsedFor_PaymentCollector:
		fallthrough // nolint
	case basetypes.AccountUsedFor_UserBenefitCold:
		fallthrough // nolint
	case basetypes.AccountUsedFor_PlatformBenefitCold:
		if h.Address == nil {
			return fmt.Errorf("invalid address")
		}
	}
	if h.Address == nil {
		return nil
	}
	if h.CoinTypeID == nil {
		return fmt.Errorf("invalid cointypeid")
	}

	info, err := pltfmwcli.GetAccountOnly(
		ctx,
		&pltfmwpb.Conds{
			CoinTypeID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.CoinTypeID},
			Address:    &basetypes.StringVal{Op: cruder.EQ, Value: *h.Address},
		},
	)
	if err != nil {
		return err
	}
	if info == nil {
		return fmt.Errorf("address not exist")
	}
	if info.UsedFor != *h.UsedFor {
		return fmt.Errorf("mismatch account")
	}

	h.ID = &info.ID
	return nil
}

func (h *createHandler) checkBackup(ctx context.Context) error {
	exist, err := pltfmwcli.ExistAccountConds(ctx, &pltfmwpb.Conds{
		CoinTypeID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.CoinTypeID},
		UsedFor:    &basetypes.Uint32Val{Op: cruder.EQ, Value: uint32(*h.UsedFor)},
		Backup:     &basetypes.BoolVal{Op: cruder.EQ, Value: false},
	})
	if err != nil {
		return err
	}
	h.backup = exist
	return nil
}

func (h *createHandler) getCoinName(ctx context.Context) error {
	if h.CoinTypeID == nil {
		return fmt.Errorf("invalid cointypeid")
	}

	coin, err := coinmwcli.GetCoin(ctx, *h.CoinTypeID)
	if err != nil {
		return err
	}
	if coin == nil {
		return fmt.Errorf("invalid coin")
	}

	h.coinName = &coin.Name
	h.checkAddressBalance = coin.CheckNewAddressBalance

	return nil
}

func (h *createHandler) createAddress(ctx context.Context) error {
	if h.coinName == nil {
		return fmt.Errorf("invalid coinname")
	}

	if h.Address == nil {
		acc, err := sphinxproxycli.CreateAddress(ctx, *h.coinName)
		if err != nil {
			return err
		}
		if acc == nil {
			return fmt.Errorf("fail create address")
		}
		h.Address = &acc.Address
	}

	if !h.checkAddressBalance {
		return nil
	}

	bal, err := sphinxproxycli.GetBalance(ctx, &sphinxproxypb.GetBalanceRequest{
		Name:    *h.coinName,
		Address: *h.Address,
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
	acc, err := pltfmwcli.CreateAccount(ctx, &pltfmwpb.AccountReq{
		CoinTypeID: h.CoinTypeID,
		UsedFor:    h.UsedFor,
		Address:    h.Address,
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

	if err := handler.checkAddress(ctx); err != nil {
		return nil, err
	}
	if h.ID != nil {
		return h.GetAccount(ctx)
	}
	if err := handler.checkBackup(ctx); err != nil {
		return nil, err
	}
	if err := handler.getCoinName(ctx); err != nil {
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
