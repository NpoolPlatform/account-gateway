package platform

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	pltfmwcli "github.com/NpoolPlatform/account-middleware/pkg/client/platform"
	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/platform"
	pltfmwpb "github.com/NpoolPlatform/message/npool/account/mw/v1/platform"
)

func (h *Handler) UpdateAccount(ctx context.Context) (*npool.Account, error) {
	info, err := pltfmwcli.GetAccountOnly(ctx, &pltfmwpb.Conds{
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
