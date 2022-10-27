package user

import (
	"context"
	"fmt"

	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/user"

	useraccmwcli "github.com/NpoolPlatform/account-middleware/pkg/client/user"
)

func DeleteAccount(ctx context.Context, id string) (*npool.Account, error) {
	info, err := GetAccount(ctx, id)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf("invalid account")
	}

	_, err = useraccmwcli.DeleteAccount(ctx, id)
	if err != nil {
		return nil, err
	}

	return info, nil
}
