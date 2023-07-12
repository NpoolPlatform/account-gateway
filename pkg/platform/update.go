package platform

import (
	"context"
	"fmt"

	pltfmwcli "github.com/NpoolPlatform/account-middleware/pkg/client/platform"
	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/platform"
	pltfmwpb "github.com/NpoolPlatform/message/npool/account/mw/v1/platform"
)

func (h *Handler) UpdateAccount(ctx context.Context) (*npool.Account, error) {
	if h.ID == nil {
		return nil, fmt.Errorf("invalid id")
	}

	info, err := pltfmwcli.GetAccount(ctx, *h.ID)
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
	boolTrue := true

	if h.Blocked != nil && *h.Blocked {
		h.Active = &boolFalse
		h.Backup = &boolTrue
	}
	if h.Active != nil && *h.Active {
		h.Blocked = &boolFalse
	}

	if _, err := pltfmwcli.UpdateAccount(ctx, &pltfmwpb.AccountReq{
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
