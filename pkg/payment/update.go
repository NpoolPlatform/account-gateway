package payment

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/payment"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	paymentmwcli "github.com/NpoolPlatform/account-middleware/pkg/client/payment"
	paymentmwpb "github.com/NpoolPlatform/message/npool/account/mw/v1/payment"
)

func (h *Handler) UpdateAccount(ctx context.Context) (*npool.Account, error) {
	info, err := paymentmwcli.GetAccountOnly(ctx, &paymentmwpb.Conds{
		ID:    &basetypes.Uint32Val{Op: cruder.EQ, Value: *h.ID},
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.EntID},
	})
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf("invalid account")
	}

	fmt.Println("info: ", info)

	if info.Blocked && (h.Blocked == nil || *h.Blocked) {
		return nil, fmt.Errorf("permission denied")
	}
	boolFalse := false

	if h.Blocked != nil && *h.Blocked {
		h.Active = &boolFalse
	}
	if h.Active != nil && *h.Active {
		h.Blocked = &boolFalse
	}

	if _, err := paymentmwcli.UpdateAccount(ctx, &paymentmwpb.AccountReq{
		ID:      h.ID,
		Active:  h.Active,
		Blocked: h.Blocked,
	}); err != nil {
		return nil, err
	}

	return h.GetAccount(ctx)
}
