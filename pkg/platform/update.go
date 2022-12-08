package platform

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/platform"

	pltfmwcli "github.com/NpoolPlatform/account-middleware/pkg/client/platform"
	pltfmwpb "github.com/NpoolPlatform/message/npool/account/mw/v1/platform"
)

func UpdateAccount(
	ctx context.Context,
	id string,
	backup, active, blocked, locked *bool,
) (
	*npool.Account, error,
) {
	acc, err := pltfmwcli.UpdateAccount(ctx, &pltfmwpb.AccountReq{
		ID:      &id,
		Active:  active,
		Blocked: blocked,
		Locked:  locked,
		Backup:  backup,
	})
	if err != nil {
		return nil, err
	}

	return GetAccount(ctx, acc.ID)
}
