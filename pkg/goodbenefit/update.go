package goodbenefit

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/goodbenefit"

	gbmwcli "github.com/NpoolPlatform/account-middleware/pkg/client/goodbenefit"
	gbmwpb "github.com/NpoolPlatform/message/npool/account/mw/v1/goodbenefit"
)

func UpdateAccount(ctx context.Context, id string, backup bool) (*npool.Account, error) {
	acc, err := gbmwcli.UpdateAccount(ctx, &gbmwpb.AccountReq{
		ID:     &id,
		Backup: &backup,
	})
	if err != nil {
		return nil, err
	}

	return GetAccount(ctx, acc.ID)
}
