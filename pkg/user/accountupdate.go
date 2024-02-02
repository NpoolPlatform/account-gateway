package user

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/user"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	useraccmwcli "github.com/NpoolPlatform/account-middleware/pkg/client/user"
	useraccmwpb "github.com/NpoolPlatform/message/npool/account/mw/v1/user"
)

func (h *Handler) UpdateAccount(ctx context.Context) (*npool.Account, error) {
	info, err := useraccmwcli.GetAccountOnly(ctx, &useraccmwpb.Conds{
		ID:    &basetypes.Uint32Val{Op: cruder.EQ, Value: *h.ID},
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.EntID},
	})
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf("invalid account")
	}
	if info.AppID != *h.AppID || info.UserID != *h.UserID {
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

	_, err = useraccmwcli.UpdateAccount(ctx, &useraccmwpb.AccountReq{
		ID:      h.ID,
		Active:  h.Active,
		Blocked: h.Blocked,
		Labels:  h.Labels,
		Memo:    h.Memo,
	})
	if err != nil {
		return nil, err
	}

	return h.GetAccount(ctx)
}
