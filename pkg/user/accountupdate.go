package user

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/user"

	useraccmwcli "github.com/NpoolPlatform/account-middleware/pkg/client/user"
	useraccmwpb "github.com/NpoolPlatform/message/npool/account/mw/v1/user"
)

func UpdateAccount(ctx context.Context, id string, active, blocked *bool, labels []string) (*npool.Account, error) {
	_, err := useraccmwcli.UpdateAccount(ctx, &useraccmwpb.AccountReq{
		ID:      &id,
		Active:  active,
		Blocked: blocked,
		Labels:  labels,
	})
	if err != nil {
		return nil, err
	}
	return GetAccount(ctx, id)
}
