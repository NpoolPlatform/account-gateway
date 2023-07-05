package user

import (
	"context"
	"fmt"

	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/user"

	useraccmwcli "github.com/NpoolPlatform/account-middleware/pkg/client/user"
	useraccmwpb "github.com/NpoolPlatform/message/npool/account/mw/v1/user"
)

func (h *Handler) UpdateAccount(ctx context.Context) (*npool.Account, error) {
	if h.ID == nil {
		return nil, fmt.Errorf("invalid id")
	}
	if h.AppID == nil {
		return nil, fmt.Errorf("invalid appID")
	}
	if h.UserID == nil {
		return nil, fmt.Errorf("invalid userID")
	}

	info, err := useraccmwcli.GetAccount(ctx, *h.ID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf("invalid account")
	}
	if info.AppID != *h.AppID || info.UserID != *h.UserID {
		return nil, fmt.Errorf("permission denied")
	}

	_, err = useraccmwcli.UpdateAccount(ctx, &useraccmwpb.AccountReq{
		ID:     h.ID,
		Labels: *h.Labels,
		Memo:   h.Memo,
	})
	if err != nil {
		return nil, err
	}

	return h.GetAccount(ctx)
}

func (h *Handler) UpdateAppUserAccount(ctx context.Context) (*npool.Account, error) { //nolint
	if h.ID == nil {
		return nil, fmt.Errorf("invalid id")
	}
	if h.AppID == nil {
		return nil, fmt.Errorf("invalid appID")
	}
	if h.UserID == nil {
		return nil, fmt.Errorf("invalid userID")
	}

	info, err := useraccmwcli.GetAccount(ctx, *h.ID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf("invalid account")
	}
	if info.AppID != *h.AppID || info.UserID != *h.UserID {
		return nil, fmt.Errorf("permission denied")
	}

	if info.Blocked && (h.Blocked == nil || *h.Blocked) {
		return nil, fmt.Errorf("permission denied")
	}
	boolFalse := false
	boolTrue := true

	if h.Blocked != nil && !*h.Blocked {
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
	})
	if err != nil {
		return nil, err
	}

	return h.GetAccount(ctx)
}
