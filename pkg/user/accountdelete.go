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

func (h *Handler) DeleteAccount(ctx context.Context) (*npool.Account, error) {
	if h.ID == nil {
		return nil, fmt.Errorf("invalid id")
	}
	if h.AppID == nil {
		return nil, fmt.Errorf("invalid appid")
	}
	if h.UserID == nil {
		return nil, fmt.Errorf("invalid userid")
	}

	info, err := useraccmwcli.GetAccount(ctx, *h.ID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf("invalid account")
	}

	exist, err := useraccmwcli.ExistAccountConds(ctx, &useraccmwpb.Conds{
		ID: &basetypes.StringVal{
			Op:    cruder.EQ,
			Value: *h.ID,
		},
		AppID: &basetypes.StringVal{
			Op:    cruder.EQ,
			Value: *h.AppID,
		},
		UserID: &basetypes.StringVal{
			Op:    cruder.EQ,
			Value: *h.UserID,
		},
	})
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("invalid account")
	}

	_, err = useraccmwcli.DeleteAccount(ctx, &useraccmwpb.AccountReq{
		ID: h.ID,
	})
	if err != nil {
		return nil, err
	}

	return h.GetAccount(ctx)
}
