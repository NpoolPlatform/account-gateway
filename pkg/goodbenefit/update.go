package goodbenefit

import (
	"context"
	"fmt"

	gbmwcli "github.com/NpoolPlatform/account-middleware/pkg/client/goodbenefit"
	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/goodbenefit"
	gbmwpb "github.com/NpoolPlatform/message/npool/account/mw/v1/goodbenefit"
)

func (h *Handler) UpdateAccount(ctx context.Context) (*npool.Account, error) {
	if h.ID == nil {
		return nil, fmt.Errorf("invalid id")
	}
	if _, err := gbmwcli.UpdateAccount(ctx, &gbmwpb.AccountReq{
		ID:      h.ID,
		Backup:  h.Backup,
		Active:  h.Active,
		Blocked: h.Blocked,
		Locked:  h.Locked,
	}); err != nil {
		return nil, err
	}

	return h.GetAccount(ctx)
}
