package payment

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/payment"

	paymentmwcli "github.com/NpoolPlatform/account-middleware/pkg/client/payment"
	paymentmwpb "github.com/NpoolPlatform/message/npool/account/mw/v1/payment"
)

func UpdateAccount(
	ctx context.Context,
	id string,
	active, blocked, locked *bool,
) (
	*npool.Account, error,
) {
	acc, err := paymentmwcli.UpdateAccount(ctx, &paymentmwpb.AccountReq{
		ID:      &id,
		Active:  active,
		Blocked: blocked,
		Locked:  locked,
	})
	if err != nil {
		return nil, err
	}

	return GetAccount(ctx, acc.ID)
}
