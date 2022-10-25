package platform

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/platform"

	gbmwcli "github.com/NpoolPlatform/account-middleware/pkg/client/platform"
	gbmwpb "github.com/NpoolPlatform/message/npool/account/mw/v1/platform"
)

func UpdateAccount(
	ctx context.Context,
	id string,
	backup, active, blocked, locked *bool,
) (
	*npool.Account, error,
) {
	acc, err := gbmwcli.UpdateAccount(ctx, &gbmwpb.AccountReq{
		ID:      &id,
		Backup:  backup,
		Active:  active,
		Blocked: blocked,
		Locked:  locked,
	})
	if err != nil {
		return nil, err
	}

	return GetAccount(ctx, acc.ID)
}
