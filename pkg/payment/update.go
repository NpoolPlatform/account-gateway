package payment

import (
	"context"
	"fmt"

	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/payment"

	paymentmwcli "github.com/NpoolPlatform/account-middleware/pkg/client/payment"
	paymentmwpb "github.com/NpoolPlatform/message/npool/account/mw/v1/payment"
)

func (h *Handler) UpdateAccount(ctx context.Context) (*npool.Account, error) {
	if h.ID == nil {
		return nil, fmt.Errorf("invalid id")
	}

	info, err := paymentmwcli.GetAccount(ctx, *h.ID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf("invalid account")
	}

	if info.Blocked && (h.Blocked == nil || *h.Blocked) {
		return nil, fmt.Errorf("permission denied")
	}

	boolFalse := false

	if h.Blocked != nil && !*h.Blocked {
		h.Active = &boolFalse
	}
	if h.Active != nil && *h.Active {
		h.Blocked = &boolFalse
	}

	if _, err := paymentmwcli.UpdateAccount(ctx, &paymentmwpb.AccountReq{
		ID:      h.ID,
		Active:  h.Active,
		Blocked: h.Blocked,
		Locked:  h.Locked,
	}); err != nil {
		return nil, err
	}

	return h.GetAccount(ctx)
}
